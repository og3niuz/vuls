package scan

import (
	"github.com/future-architect/vuls/config"
	"github.com/future-architect/vuls/models"
	"github.com/future-architect/vuls/util"
)

// inherit OsTypeInterface
type oracle struct {
	redhatBase
}

// NewAmazon is constructor
func newOracle(c config.ServerInfo) *oracle {
	r := &oracle{
		redhatBase{
			base: base{
				osPackages: osPackages{
					Packages:  models.Packages{},
					VulnInfos: models.VulnInfos{},
				},
			},
		},
	}
	r.log = util.NewCustomLogger(c)
	r.setServerInfo(c)
	return r
}

func (o *oracle) checkDeps() error {
	if config.Conf.Fast {
		return o.execCheckDeps(o.depsFast())
	} else if config.Conf.FastRoot {
		return o.execCheckDeps(o.depsFastRoot())
	} else {
		return o.execCheckDeps(o.depsDeep())
	}
}

func (o *oracle) depsFastRoot() []string {
	if config.Conf.Offline {
		//TODO
		// return []string{"yum-plugin-ps"}
	}

	majorVersion, _ := o.Distro.MajorVersion()
	switch majorVersion {
	case 5:
		return []string{
			"yum-utils",
			"yum-security",
		}
	case 6:
		return []string{
			"yum-utils",
			"yum-plugin-security",
			//TODO
			// return []string{"yum-plugin-ps"}
		}
	default:
		return []string{
			"yum-utils",
			//TODO
			// return []string{"yum-plugin-ps"}
		}
	}
}

func (o *oracle) depsDeep() []string {
	majorVersion, _ := o.Distro.MajorVersion()
	switch majorVersion {
	case 5:
		return []string{
			"yum-utils",
			"yum-security",
			"yum-changelog",
		}
	case 6:
		return []string{
			"yum-utils",
			"yum-plugin-security",
			"yum-plugin-changelog",
			//TODO
			// return []string{"yum-plugin-ps"}
		}
	default:
		return []string{
			"yum-utils",
			"yum-plugin-changelog",
			//TODO
			// return []string{"yum-plugin-ps"}
		}
	}
}

func (o *oracle) checkIfSudoNoPasswd() error {
	if config.Conf.Fast {
		return o.execCheckIfSudoNoPasswd(o.nosudoCmdsFast())
	} else if config.Conf.FastRoot {
		return o.execCheckIfSudoNoPasswd(o.nosudoCmdsFastRoot())
	} else {
		return o.execCheckIfSudoNoPasswd(o.nosudoCmdsDeep())
	}
}

func (o *oracle) nosudoCmdsFast() []cmd {
	return []cmd{}
}

func (o *oracle) nosudoCmdsFastRoot() []cmd {
	cmds := []cmd{{"needs-restarting", exitStatusZero}}
	if config.Conf.Offline {
		return cmds
	}

	majorVersion, _ := o.Distro.MajorVersion()
	if majorVersion < 6 {
		return []cmd{
			{"yum --color=never repolist", exitStatusZero},
			{"yum --color=never list-security --security", exitStatusZero},
			{"yum --color=never info-security", exitStatusZero},
			{"repoquery -h", exitStatusZero},
		}
	}
	return append(cmds,
		cmd{"yum --color=never repolist", exitStatusZero},
		cmd{"yum --color=never --security updateinfo list updates", exitStatusZero},
		cmd{"yum --color=never --security updateinfo updates", exitStatusZero},
		cmd{"repoquery -h", exitStatusZero})
}

func (o *oracle) nosudoCmdsDeep() []cmd {
	return append(o.nosudoCmdsFastRoot(),
		cmd{"yum --color=never repolist", exitStatusZero},
		cmd{"yum changelog all updates", exitStatusZero})
}
