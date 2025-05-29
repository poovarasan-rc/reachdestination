package accountholder

import (
	"appsec/db"
	"appsec/model"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Response struct
type ResponseStruct struct {
	model.AccountStruct
	model.Debug
}

func AccountHolder(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	log.Println("AccountHolder(+)")
	var lReq model.AccountStruct
	var lRespRec ResponseStruct

	lRespRec.Sts = "S"
	lRespRec.Msg = ""

	// Reading request body
	lBody, lErr := ioutil.ReadAll(r.Body)
	if lErr != nil {
		lRespRec.Sts = "E"
		lRespRec.Msg = "AACH01 : " + lErr.Error()
		log.Println(lRespRec)
		goto Marshal
	}
	// Unmarshall the body into respective struct
	lErr = json.Unmarshal(lBody, &lReq)
	if lErr != nil {
		lRespRec.Sts = "E"
		lRespRec.Msg = "AACH02 : " + lErr.Error()
		log.Println(lRespRec)
		goto Marshal
	}
	// Call the method to perform the action
	lRespRec.AccountStruct, lErr = Perf_Acc_Action(lReq, r.Method)
	if lErr != nil {
		lRespRec.Sts = "E"
		lRespRec.Msg = "AACH03 : " + lErr.Error()
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
	log.Println("AccountHolder(-)")

}

// Method to store the Account
func Perf_Acc_Action(pRec model.AccountStruct, pMethod string) (model.AccountStruct, error) {
	var lErr error
	var lResp model.AccountStruct
	switch pMethod {
	case http.MethodGet:
		// Get Account details
		lResp, lErr = GetAccount(pRec.AccID)
		if lErr != nil {
			return lResp, lErr
		}
	case http.MethodPost:
		// Insert Account details
		lErr = InsertAccount(pRec)
		if lErr != nil {
			return lResp, lErr
		}
	case http.MethodPut:
		// Update Account details
		lErr = ModifyAccount(pRec)
		if lErr != nil {
			return lResp, lErr
		}
	case http.MethodDelete:
		// Delete Account details
		lErr = DeleteAccount(pRec.AccID)
		if lErr != nil {
			return lResp, lErr
		}
	default:
		return lResp, fmt.Errorf("Invalid HTTP method - " + pMethod)
	}
	return lResp, nil
}

// Method to Get Account details
func GetAccount(pAccId string) (model.AccountStruct, error) {
	var lResp model.AccountStruct

	query := `select accname,email from account_holders ah where accid = ?`
	lErr := db.GDBCon.QueryRow(query, pAccId).Scan(&lResp.AccName, &lResp.Email)
	if lErr != nil {
		return lResp, fmt.Errorf("Account not available - " + pAccId)
	}
	return lResp, nil
}

// Method to Insert Account details
func InsertAccount(pRec model.AccountStruct) error {

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString([]byte(pRec.AccID + pRec.AccName + pRec.Email))

	lQuery := `INSERT into account_holders (accid,accname,email,appsecret,createdby,createddate)
		values(?,?,?,?,?,Now())`
	_, lErr := db.GDBCon.Exec(lQuery, pRec.AccID, pRec.AccName, pRec.Email, encoded, pRec.Email)
	if lErr != nil {
		return lErr
	}
	return nil
}

// Method to Update Account details
func ModifyAccount(pRec model.AccountStruct) error {
	lQuery := `UPDATE appscrt.account_holders
		SET  email = ?, accname=?, updatedby=?, updateddate=Now()
		WHERE accid = ?;`
	lResults, lErr := db.GDBCon.Exec(lQuery, pRec.Email, pRec.AccName, pRec.Email, pRec.AccID)
	if lErr != nil {
		return lErr
	}
	lResult, _ := lResults.RowsAffected()
	if lResult == 0 {
		return fmt.Errorf("Unable to update, Please provide valide Account ID - " + pRec.AccID)
	}
	return nil
}

// Method to Delete Account details
func DeleteAccount(pAccID string) error {

	lTrans, lErr := db.GDBCon.Begin()
	if lErr != nil {
		return lErr
	}

	// Delete the destinations of the account
	lQuery1 := `DELETE from destination 
		where accid = (select id from account_holders ah where accid = ?)`
	_, lErr = lTrans.Exec(lQuery1, pAccID)
	if lErr != nil {
		lTrans.Rollback()
		return lErr
	}

	lQuery := `delete from account_holders 
		WHERE accid = ?;`
	lResults, lErr := lTrans.Exec(lQuery, pAccID)
	if lErr != nil {
		lTrans.Rollback()
		return lErr
	}

	lTrans.Commit()

	lResult, _ := lResults.RowsAffected()
	if lResult == 0 {
		return fmt.Errorf("Unable to remove, Please provide valide Account ID - " + pAccID)
	}
	return nil
}
