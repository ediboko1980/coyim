package gui

import (
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/session"
)

//TODO: once glib signals are removed from account, it could be a sessionManager
// and this could be in a client package (ui agnostic)
type accountManager struct {
	accounts          []*account
	events            chan interface{}
	saveConfiguration func()
}

func newAccountManager() *accountManager {
	return &accountManager{
		events: make(chan interface{}, 10),
	}
}

func (m *accountManager) addAccount(appConfig *config.Accounts, account *config.Account) {
	acc, err := newAccount(appConfig, account)
	if err != nil {
		//TODO error
		return
	}

	acc.session.SaveConfiguration = m.saveConfiguration
	acc.session.Subscribe(m.events)
	m.accounts = append(m.accounts, acc)
}

func (m *accountManager) buildAccounts(appConfig *config.Accounts) {
	m.accounts = make([]*account, 0, len(appConfig.Accounts))
	hasConfUpdates := false
	for _, accountConf := range appConfig.Accounts {
		hasUpdate, err := accountConf.EnsurePrivateKey()
		if err != nil {
			continue
		}
		hasConfUpdates = hasConfUpdates || hasUpdate
		m.addAccount(appConfig, accountConf)
	}
	if hasConfUpdates {
		m.saveConfiguration()
	}
}

func (m *accountManager) findAccountForSession(s *session.Session) *account {
	for _, a := range m.accounts {
		if a.session == s {
			return a
		}
	}

	return nil
}

func (m *accountManager) findAccountForUsername(s string) *account {
	for _, a := range m.accounts {
		if a.session.CurrentAccount.Is(s) {
			return a
		}
	}

	return nil
}

func (m *accountManager) addNewAccountsFromConfig(appConfig *config.Accounts) {
	for _, configAccount := range appConfig.Accounts {
		var found bool
		for _, acc := range m.accounts {
			if acc.session.CurrentAccount.ID() == configAccount.ID() {
				found = true
				break
			}
		}

		if found {
			continue
		}

		m.addAccount(appConfig, configAccount)
	}
}
