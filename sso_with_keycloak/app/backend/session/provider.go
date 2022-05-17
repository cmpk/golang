//
// お世話になったサイト：
//    6.2 Goはどのようにしてsessionを使用するか
//      https://astaxie.gitbooks.io/build-web-application-with-golang/content/ja/06.2.html
//

package session

import (
	"fmt"
)

type Provider struct {
	sessions map[string]Session
}

func CreateProvider() Provider {
	return Provider{sessions: map[string]Session{}}
}

func (provider *Provider) createEmptySession() Session {
	return CreateSession("")
}

func (provider *Provider) SessionInit(sid string) (Session, error) {
	var session Session
	session, ok := provider.sessions[sid]
	if ok {
		return provider.createEmptySession(), fmt.Errorf("already exist. sessionId: %s", sid)
	}
	session = CreateSession(sid)
	provider.sessions[sid] = session
	return session, nil
}

func (provider *Provider) SessionRead(sid string) (Session, error) {
	var session Session
	session, ok := provider.sessions[sid]
	if !ok {
		return provider.createEmptySession(), fmt.Errorf("not found session ID. sessionId: %s", sid)
	}
	return session, nil
}

func (provider *Provider) SessionDestroy(sid string) error {
	delete(provider.sessions, sid)
	return nil
}

func (provider *Provider) SessionGC(maxLifeTime int64) {
	panic("Not Implemented")
}
