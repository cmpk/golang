//
// お世話になったサイト：
//    6.2 Goはどのようにしてsessionを使用するか
//      https://astaxie.gitbooks.io/build-web-application-with-golang/content/ja/06.2.html
//

package session

import (
	"fmt"
)

type Session struct {
	sessionId string
	data      map[string]string
}

func CreateSession(sessionId string) Session {
	return Session{sessionId: sessionId, data: map[string]string{}}
}

func (session *Session) Set(key string, value string) error {
	session.data[key] = value
	return nil
}

func (session *Session) Get(key string) (string, error) {
	v, ok := session.data[key]
	if !ok {
		return "", fmt.Errorf("not found key. sessionId: %s, key: %s", session.sessionId, key)
	}
	return v, nil
}

func (session *Session) Delete(key string) error {
	_, ok := session.data[key]
	if !ok {
		return fmt.Errorf("cannot delete key. sessionId: %s, key: %s", session.sessionId, key)
	}
	delete(session.data, key)
	return nil
}

func (session *Session) SessionID() string {
	return session.sessionId
}
