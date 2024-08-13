package cmd

import (
	"log"

	entity "upsun.com/lib-upsun/entity"
)

func DisplayUnsupportedIntegrations(projectContext entity.ProjectGlobal) {
	log.Printf("upsun_clone do not clone Integrations. Please add them manually !")
}
