package cmd

import (
	"log"

	entity "github.com/upsun/lib-sun/entity"
)

func DisplayUnsupportedIntegrations(projectContext entity.ProjectGlobal) {
	log.Printf("upsun_clone do not clone Integrations. Please add them manually !")
}
