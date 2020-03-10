package hwcallback

// //华为云短信回调
// func postNoteback(w http.ResponseWriter, r *http.Request) {
// 	v := new(HWYNote)
// 	jData := ioutil.ReadAll(r.Body)
// 	if err := json.Unmarshal(jData, &v); err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(string(jData))
// 	w.WriteHeader(http.StatusOK)
// }
