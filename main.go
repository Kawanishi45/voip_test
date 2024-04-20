package main

import (
  "fmt"
  "github.com/google/uuid"
  "github.com/sideshow/apns2"
  "github.com/sideshow/apns2/certificate"
  "log"
)

func main() {
  // .p12形式の証明書ファイルのパスとパスワードを指定
  cert, err := certificate.FromP12File("./cert.p12", certFilePass)
  if err != nil {
    log.Fatal("Cert error:", err)
  }

  // APNsクライアントの作成
  client := apns2.NewClient(cert).Development() // Development or Production

  // VoIP通知を送るための通知オブジェクトの作成
  notification := &apns2.Notification{}
  notification.DeviceToken = iphoneDeviceToken
  notification.Topic = notificationTopic
  //notification.Payload = []byte(`{"aps":{"alert":"Hello VoIP Push","sound":"default"}}`) // JSON形式のペイロード
  var u uuid.UUID
  u, err = uuid.NewRandom()
  if err != nil {
    log.Fatal("Error:", err)
  }
  payload := fmt.Sprintf(`{"aps":{"alert":{"uuid":"%s","incoming_caller_id":"1","incoming_caller_name":"user_name"}}}`, u.String())
  log.Println(payload)
  notification.Payload = []byte(payload) // JSON形式のペイロード
  // {
  //    "aps": {
  //        "alert": {
  //          "uuid": <Version 4 UUID (e.g.: https://www.uuidgenerator.net/version4) >,
  //          "incoming_caller_id": <your service user id>,
  //          "incoming_caller_name": <your service user name>,
  //        }
  //    }
  //}
  //notification.Payload = payload.NewPayload().Alert("hello").Badge(1).Custom("key", "val")
  notification.PushType = apns2.PushTypeVOIP // VoIP通知を送るためのPushType

  // 通知を送信
  var res *apns2.Response
  res, err = client.Push(notification)
  if err != nil {
    log.Fatal("Error:", err)
  }

  // 送信結果の確認
  log.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}
