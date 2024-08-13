package cmd

import (
	"encoding/json"
	"log"

	entity "upsun.com/lib-upsun/entity"
	utils "upsun.com/lib-upsun/utility"
)

type getAccess struct {
	Items []entity.ProjectAccess
}

func UsersRead(projectContext entity.ProjectGlobal) {
	log.Print("Read Users...")

	payload := []string{"-X", "GET", "/user-access"}
	jsonContent, _ := utils.CallCLI(projectContext, "project:curl", payload...)

	var access getAccess
	if err := json.Unmarshal(jsonContent, &access); err != nil {
		log.Printf("failed to unmarshal response: %s", err)
	}

	// TODO : resync array (remove not use)
	for _, acces := range access.Items {
		log.Printf("Find Access: %q", acces.UserId)
		projectContext.Access[acces.UserId] = acces
	}
}

func UsersWrite(projectContext entity.ProjectGlobal) {
	log.Print("Write Users...")

	var result utils.PshResult
	payloadBase := []string{"-X", "POST", "/user-access", "--json"}

	for _, acces := range projectContext.Access {
		log.Printf("Write user access: %q", acces.UserId)

		_, ok := projectContext.Users[acces.UserId]
		if ok {
			log.Print("It is me ! nothing to do")
		} else {
			//TODO need to add a Hack !!!

			// DTO (dynamic)
			var dtos []entity.ProjectAccess
			dto := acces                    // Make a copy (for not invalidate original)
			dto.AddAuto = &[]bool{true}[0]  // Disable auto add (hack)
			dtos = append(dtos, dto)        // switch to array
			output, _ := json.Marshal(dtos) // Convert to JSON (by Marshal)
			dtoJson := string(output)       // Transform to string

			// CREATE case.
			payloadInsert := append(payloadBase, dtoJson)
			result = utils.CallAPI(projectContext, payloadInsert...)

			// UPDATE case.
			if result.Title == "Conflict" {
				dto.UserId = ""
				dto.AddAuto = nil
				output, _ := json.Marshal(dto) // Convert to JSON (by Marshal)
				dtoJson := string(output)      // Transform to string

				payloadUpdate := []string{"-X", "PATCH",
					"/user-access/" + acces.UserId,
					"--json", dtoJson}
				utils.CallAPI(projectContext, payloadUpdate...)
			}
		}

	}
}
