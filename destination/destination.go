package destination

import (
	"appsec/db"
	"appsec/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Response struct
type Response struct {
	model.Destination
	model.Debug
}

func Designation(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	log.Println("Designation(+)")

	var lRespRec Response
	var lDestRec model.Destination

	lRespRec.Sts = "S"
	lRespRec.Msg = ""

	// Reading request body
	lBody, lErr := ioutil.ReadAll(r.Body)
	if lErr != nil {
		lRespRec.Sts = "E"
		lRespRec.Msg = "DDT01 : " + lErr.Error()
		log.Println(lRespRec)
		goto Marshal
	}
	// Unmarshall the body into respective struct
	lErr = json.Unmarshal(lBody, &lDestRec)
	if lErr != nil {
		lRespRec.Sts = "E"
		lRespRec.Msg = "DDT02 : " + lErr.Error()
		log.Println(lRespRec)
		goto Marshal
	}
	// Call the method to perform the action
	lRespRec.Destination, lErr = Perf_Dest_Action(lDestRec, r.Method)
	if lErr != nil {
		lRespRec.Sts = "E"
		lRespRec.Msg = "DDT03 : " + lErr.Error()
		log.Println(lRespRec)
		goto Marshal
	}
Marshal:
	data, err := json.Marshal(lRespRec)
	if err != nil {
		fmt.Fprintf(w, "Error taking data"+err.Error())
	} else {
		fmt.Fprintf(w, string(data))
	}
	log.Println("Designation(-)")

}

// Method to store the destinations
func Perf_Dest_Action(pRec model.Destination, pMethod string) (model.Destination, error) {

	var lDesRec model.Destination
	var lErr error

	switch pMethod {
	case http.MethodGet:
		// Get destination details
		lDesRec, lErr = GetDest(pRec.DesID)
		if lErr != nil {
			return lDesRec, lErr
		}
	case http.MethodPost:
		// Insert destination details
		lErr = InsertDest(pRec)
		if lErr != nil {
			return lDesRec, lErr
		}
	case http.MethodPut:
		// Update destination details
		lErr = ModifyDest(pRec)
		if lErr != nil {
			return lDesRec, lErr
		}
	case http.MethodDelete:
		// Delete destination details
		lErr = DeleteDest(pRec.DesID)
		if lErr != nil {
			return lDesRec, lErr
		}
	default:
		return lDesRec, fmt.Errorf("Invalid HTTP method - " + pMethod)
	}
	return lDesRec, nil
}

// Method to Get destination details
func GetDest(pDesc string) (model.Destination, error) {
	var lHeaders string
	var lDesRec model.Destination

	query := `SELECT d.desId,ah.accid,d.url,d.method,d.headers
		from account_holders ah ,destination d 
		where ah.id = d.accid
		and d.desid = ?`
	lErr := db.GDBCon.QueryRow(query, pDesc).Scan(&lDesRec.DesID, &lDesRec.AccountID, &lDesRec.URL, &lDesRec.HTTPMethod, &lHeaders)
	if lErr != nil {
		return lDesRec, lErr
	}
	if lHeaders != "" {
		err := json.Unmarshal([]byte(lHeaders), &lDesRec.Headers)
		if err != nil {
			fmt.Println("Error unmarshalling:", err)
			return lDesRec, err
		}
	}
	return lDesRec, nil
}

// Method to Insert destination details
func InsertDest(pRec model.Destination) error {
	Id, _ := getAccountId(pRec.AccountID)
	if Id == 0 {
		return fmt.Errorf("Account not available - " + pRec.AccountID)
	}
	data2, lErr := json.Marshal(pRec.Headers)
	if lErr != nil {
		return lErr
	}
	lQuery := `INSERT INTO destination (desid,accid, url, method, headers, createdby, createddate )
		-- VALUES(?, ?, ?, ?, ?, Now());
		SELECT ?,?,?,?,?,?,Now()
		where not exists (select 1 from destination d 
						where URL = ? and accid = ?)`
	lResults, lErr := db.GDBCon.Exec(lQuery, pRec.DesID, Id, pRec.URL, pRec.HTTPMethod, string(data2), "AUTOBOT", pRec.URL, Id)
	if lErr != nil {
		return lErr
	}
	lResult, _ := lResults.RowsAffected()
	if lResult == 0 {
		return fmt.Errorf("URL already exist for the account - " + pRec.URL)
	}
	return nil
}

// Method to Update destination details
func ModifyDest(pRec model.Destination) error {
	Id, _ := getAccountId(pRec.AccountID)
	if Id == 0 {
		return fmt.Errorf("Account not available - " + pRec.AccountID)
	}
	data2, lErr := json.Marshal(pRec.Headers)
	if lErr != nil {
		return lErr
	}
	lQuery := `UPDATE destination
		SET accid=?, url=?, method=?, headers=?, updatedby=?, updateddate=Now()
		WHERE desid = ?;`
	lResults, lErr := db.GDBCon.Exec(lQuery, Id, pRec.URL, pRec.HTTPMethod, string(data2), "AUTOBOT", pRec.DesID)
	if lErr != nil {
		return lErr
	}
	lResult, _ := lResults.RowsAffected()
	if lResult == 0 {
		return fmt.Errorf("Unable to update, Destination ID not available - " + pRec.DesID)
	}
	return nil
}

// Method to Delete destination details
func DeleteDest(pDesId string) error {

	lQuery := `delete from destination
		WHERE desid = ?;`
	lResults, lErr := db.GDBCon.Exec(lQuery, pDesId)
	if lErr != nil {
		return lErr
	}
	lResult, _ := lResults.RowsAffected()
	if lResult == 0 {
		return fmt.Errorf("Unable to delete, Destination ID not available - " + pDesId)
	}
	return nil
}

// Method to Get Account_ID for the app_secret
func getAccountId(pSecret string) (int, error) {
	var id int

	query := `select id from account_holders ah 
	where accid = ?`
	err := db.GDBCon.QueryRow(query, pSecret).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, nil
}
