package incomingdata

import (
	"appsec/apiUtil"
	"appsec/db"
	"appsec/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Response Struct
type ResponseStruct struct {
	Data []model.ConstructData `json:"data"`
	model.Debug
}

func PassIncomingData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "App_Secret, Accept, Content-Type")

	log.Println("PassIncomingData(+)")

	var lErr error
	var lAppSecret string
	var lResp ResponseStruct
	var lJsondata map[string]string

	lAppSecret = r.Header.Get("App_Secret")

	// Validate the request
	lJsondata, lErr = validateReq(r, lAppSecret)
	if lErr != nil {
		lResp.Sts = "E"
		lResp.Msg = lErr.Error()
		goto Marshal
	}

	// Retrieve the destination details of the appsecret
	lResp.Data, lErr = RetrieveDestDet(lAppSecret, lJsondata)
	if lErr != nil {
		lResp.Sts = "E"
		lResp.Msg = lErr.Error()
		goto Marshal
	}
	// Construct the req and call the destination API.
	CallDestinationApi(lResp.Data, lJsondata)

	// ✅ Valid data
	lResp.Sts = "S"
	lResp.Msg = "Data received successfully"

Marshal:
	data, lErr := json.Marshal(lResp)
	if lErr != nil {
		fmt.Fprintf(w, "Error taking data"+lErr.Error())
	} else {
		fmt.Fprintf(w, string(data))
	}

	log.Println("PassIncomingData(-)")
}

// Validate the request
func validateReq(r *http.Request, appSecret string) (map[string]string, error) {

	var jsondata map[string]string

	// 1. Only POST method allowed
	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("Only POST method allowed")
	}

	// 2. Content-Type must be JSON
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("Invalid Data")
	}

	// 3. Check for App-Secret header
	if appSecret == "" {
		return nil, fmt.Errorf("Un Authenticate")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Invalid Data")
	}

	err = json.Unmarshal(body, &jsondata)
	if err != nil {
		return nil, fmt.Errorf("Invalid Data")
	}
	return jsondata, nil
}

// Retrieve the destination details of the appsecret
func RetrieveDestDet(pAppSecret string, pJsondata map[string]string) ([]model.ConstructData, error) {
	var lArr []model.ConstructData

	query := `
        SELECT d.url,d.method,d.headers
		from account_holders ah,destination d 
		where ah.id = d.accid
		and ah.appsecret = ?;
    `
	rows, err := db.GDBCon.Query(query, pAppSecret)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var lRec model.ConstructData

		err := rows.Scan(&lRec.Url, &lRec.Method, &lRec.Headers)
		if err != nil {
			return nil, err
		}

		lArr = append(lArr, lRec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(lArr) == 0 {
		return nil, fmt.Errorf("Destionations not available for the secret (or) Invalid secret")
	}

	// log.Println("lArr - ", lArr)
	return lArr, nil
}

// Construct the req and call the destination API.
func CallDestinationApi(pDataArr []model.ConstructData, pJsondata map[string]string) {

	// Range each destination
	for idx, val := range pDataArr {

		method := strings.ToUpper(val.Method)
		// Process data based on methodstring
		switch method {
		case "GET":
			// Convert map to query string: ?key1=val1&key2=val2
			params := url.Values{}
			for k, v := range pJsondata {
				params.Add(k, v)
			}
			val.Url += "?" + params.Encode() // final full GET URL
			pDataArr[idx].Url = val.Url
		default:
			// Convert to JSON string
			jsonBytes, err := json.Marshal(pJsondata)
			if err != nil {
				log.Println(err)
			}
			val.JBody = string(jsonBytes)
			pDataArr[idx].JBody = val.JBody
		}

		// Construct the headers
		lHeaderArr := constheader(val.Headers)

		// log.Println(val.Url, "-", val.Method, "-", val.JBody, "-", lHeaderArr)
		_, lErr := apiUtil.Api_call(val.Url, val.Method, val.JBody, lHeaderArr, "")
		if lErr != nil {
			log.Println(lErr)
		} else {
			log.Println("Api called Successfully")
		}
	}
}

// Construct the headers
func constheader(pHeader string) []apiUtil.HeaderDetails {
	var lHeaderArr []apiUtil.HeaderDetails
	if pHeader != "" {
		var lHeaderMap map[string]string
		err := json.Unmarshal([]byte(pHeader), &lHeaderMap)
		if err != nil {
			fmt.Println("Error unmarshalling:", err)
		} else {
			for lkey, lval := range lHeaderMap {
				lHeaderArr = append(lHeaderArr, apiUtil.HeaderDetails{Key: lkey, Value: lval})
			}
		}
	}
	return lHeaderArr
}
