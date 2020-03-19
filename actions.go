package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/conexionBD"
	"github.com/xubiosueldos/conexionBD/Siradig/structSiradig"
	"github.com/xubiosueldos/framework"
)

type IdsAEliminar struct {
	Ids []int `json:"ids"`
}

var nombreMicroservicio string = "siradig"

func Healthy(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Healthy"))
}

func SiradigList(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)
		var siradigs []structSiradig.Siradig

		//Autocompleta la información básica de Legajo, si quiero autocompletar un substruct de Legajo (Hijos por ejemplo) se pone Legajo.Hijos
		db.Preload("Legajo").Find(&siradigs)
		framework.RespondJSON(w, http.StatusOK, siradigs)

	}

}

func SiradigShow(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)
		siradig_id := params["id"]
		p_siradigid, err := strconv.Atoi(siradig_id)
		if err != nil {
			fmt.Println(err)
		}
		framework.CheckParametroVacio(p_siradigid, w)
		var siradig structSiradig.Siradig

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)
		defer conexionBD.CerrarDB(db)

		if err := db.Set("gorm:auto_preload", true).First(&siradig, "id = ?", siradig_id).Error; gorm.IsRecordNotFoundError(err) {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, siradig)
	}

}

func SiradigAdd(w http.ResponseWriter, r *http.Request) {
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		decoder := json.NewDecoder(r.Body)

		var siradig_data structSiradig.Siradig

		if err := decoder.Decode(&siradig_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)
		defer conexionBD.CerrarDB(db)
		cero, _ := strconv.ParseFloat("0", 64)
		if canInsertUpdate(&siradig_data, db) {

			for i := 0; i < len(siradig_data.Detallecargofamiliarsiradig); i++ {
				if siradig_data.Detallecargofamiliarsiradig[i].Montoanual == nil {
					siradig_data.Detallecargofamiliarsiradig[i].Montoanual = &cero
				}
			}

			if err := db.Create(&siradig_data).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			framework.RespondJSON(w, http.StatusCreated, siradig_data)
		} else {

			framework.RespondError(w, http.StatusInternalServerError, "El Siradig que desea guardar, ya existe")
		}
	}

}

func SiradigUpdate(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)

		siradig_id, _ := strconv.ParseInt(params["id"], 10, 64)
		p_siradigid := int(siradig_id)

		if p_siradigid == 0 {
			framework.RespondError(w, http.StatusNotFound, framework.IdParametroVacio)
			return
		}

		decoder := json.NewDecoder(r.Body)

		var siradig_data structSiradig.Siradig

		if err := decoder.Decode(&siradig_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		siradigid := siradig_data.ID

		if p_siradigid == siradigid || siradigid == 0 {

			siradig_data.ID = p_siradigid

			tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
			db := conexionBD.ObtenerDB(tenant)

			defer conexionBD.CerrarDB(db)
			cero, _ := strconv.ParseFloat("0", 64)
			if canInsertUpdate(&siradig_data, db) {

				for i := 0; i < len(siradig_data.Detallecargofamiliarsiradig); i++ {
					if siradig_data.Detallecargofamiliarsiradig[i].Montoanual == nil {
						siradig_data.Detallecargofamiliarsiradig[i].Montoanual = &cero
					}
				}

				if err := db.Save(&siradig_data).Error; err != nil {
					framework.RespondError(w, http.StatusInternalServerError, err.Error())
					return
				}

				framework.RespondJSON(w, http.StatusOK, siradig_data)
			} else {
				framework.RespondError(w, http.StatusInternalServerError, "El Siradig que desea guardar, ya existe")
			}

		} else {
			framework.RespondError(w, http.StatusNotFound, framework.IdParametroDistintoStruct)
			return
		}
	}

}

func canInsertUpdate(siradig *structSiradig.Siradig, db *gorm.DB) bool {
	var idSiradig int
	caninsertupdate := true
	periodosiradiganio := siradig.Periodosiradig.Year()
	periodosiradigmes := siradig.Periodosiradig.Format("01")
	sql := "SELECT id FROM siradig WHERE legajoid = " + strconv.Itoa(*siradig.Legajoid) + " AND to_char(periodosiradig, 'MM') = '" + periodosiradigmes + "' AND to_char(periodosiradig, 'YYYY') = '" + strconv.Itoa(periodosiradiganio) + "' AND siradig.ID != " + strconv.Itoa(siradig.ID)
	db.Raw(sql).Row().Scan(&idSiradig)
	if idSiradig != 0 {
		caninsertupdate = false
	}
	return caninsertupdate
}

func SiradigRemove(w http.ResponseWriter, r *http.Request) {
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		//Para obtener los parametros por la url
		params := mux.Vars(r)
		siradig_id := params["id"]

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)
		defer conexionBD.CerrarDB(db)

		//--Borrado Fisico
		if err := db.Unscoped().Where("id = ?", siradig_id).Delete(structSiradig.Siradig{}).Error; err != nil {

			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, "El Siradig "+siradig_id+framework.MicroservicioEliminado)
	}
}

func SiradigRemoveMasivo(w http.ResponseWriter, r *http.Request) {
	var resultadoDeEliminacion = make(map[int]string)
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		var idsEliminar IdsAEliminar
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&idsEliminar); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		if len(idsEliminar.Ids) > 0 {
			for i := 0; i < len(idsEliminar.Ids); i++ {
				siradig_id := idsEliminar.Ids[i]
				if err := db.Unscoped().Where("id = ?", siradig_id).Delete(structSiradig.Siradig{}).Error; err != nil {
					//framework.RespondError(w, http.StatusInternalServerError, err.Error())
					resultadoDeEliminacion[siradig_id] = string(err.Error())

				} else {
					resultadoDeEliminacion[siradig_id] = "Fue eliminado con exito"
				}
			}
		} else {
			framework.RespondError(w, http.StatusInternalServerError, "Seleccione por lo menos un registro")
		}

		framework.RespondJSON(w, http.StatusOK, resultadoDeEliminacion)
	}

}
