package service

//var chatConnMap = sync.Map{}
//
//func RunMessageServer() {
//	listen, err := net.Listen("tcp", "127.0.0.1:9090")
//	if err != nil {
//		fmt.Printf("run server failed: %v\n", err)
//		return
//	}
//	for {
//		conn, err := listen.Accept()
//		if err != nil {
//			fmt.Printf("accept connection failed: %v\n", err)
//			continue
//		}
//		go process(conn)
//	}
//}
//
//func process(conn net.Conn) {
//	defer conn.Close()
//	var buf [256]byte
//	for {
//		n, err := conn.Read(buf[:])
//		if n == 0 {
//			if err == io.EOF {
//				break
//			}
//			fmt.Printf("read message failed: %v\n", err)
//			continue
//		}
//		var event = controller.MessageSendEvent{}
//		_ = json.Unmarshal(buf[:n], &event)
//		fmt.Printf("receive message %+v\n", event)
//		chatKey := fmt.Sprintf("%d+%d", event.FromUserId, event.ToUserId)
//		if len(event.MsgContent) == 0 { // first connection
//			chatConnMap.Store(chatKey, conn)
//			continue
//		}
//		chatKey = fmt.Sprintf("%d+%d", event.ToUserId, event.FromUserId)
//		writeConn, exist := chatConnMap.Load(chatKey)
//		if !exist {
//			fmt.Printf("user %v off line\n", event.ToUserId)
//			continue
//		}
//		pushEvent := controller.MessagePushEvent{
//			FromUserId: event.FromUserId,
//			MsgContent: event.MsgContent,
//		}
//		pushData, _ := json.Marshal(pushEvent)
//
//		_, err = writeConn.(net.Conn).Write(pushData)
//		if err != nil {
//			fmt.Printf("push message failed: %v\n", err)
//		}
//	}
//}
