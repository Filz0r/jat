package services

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/filz0r/jat/internal/database"
	"github.com/filz0r/jat/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type applicationKindMapping struct {
	kind database.ApplicationStatusKind
	name string
}

func (sm *ServiceManager) CreateInitialApplicationStatus(userID uuid.UUID) error {
	initialStatuses := []applicationKindMapping{
		{
			kind: database.Irrelevant,
			name: "Didn't Apply",
		},
		{
			kind: database.Applied,
			name: "Applied",
		},
		{
			kind: database.Ghosted,
			name: "Ghosted",
		},
		{
			kind: database.Rejected,
			name: "Rejected",
		},
		{
			kind: database.Interviewed,
			name: "Contacted",
		},
		{
			kind: database.Interviewed,
			name: "1st Interview",
		},
		{
			kind: database.Interviewed,
			name: "2nd Interview",
		},
		{
			kind: database.Interviewed,
			name: "3rd Interview",
		},
		{
			name: "Accepted",
			kind: database.Accepted,
		},
	}
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	data := make([]database.ApplicationStatus, 0, len(initialStatuses))
	for _, initialStatus := range initialStatuses {
		temp := database.ApplicationStatus{
			Status: initialStatus.name,
			UserID: userID,
			Kind:   initialStatus.kind,
		}
		data = append(data, temp)
	}
	result := sm.db.Create(&data)
	if result.Error != nil {
		return result.Error
	}
	// this is kinda garbage but it should work
	err := sm.SetUserDefaultApplicationStatus(userID, data[1].ID)
	return err
}

func (sm *ServiceManager) CreateApplicationStatus(
	userID uuid.UUID,
	status string,
	kind database.ApplicationStatusKind,
) (database.ApplicationStatus, error) {
	if sm.db == nil {
		return database.ApplicationStatus{}, errors.New("database not initialized")
	}
	applicationStatus := database.ApplicationStatus{
		Status: status,
		UserID: userID,
		Kind:   kind,
	}
	result := sm.db.Create(&applicationStatus)

	if result.Error != nil {
		return database.ApplicationStatus{}, result.Error
	}

	return applicationStatus, nil
}

func (sm *ServiceManager) UpdateApplicationStatus(
	userID uuid.UUID,
	statusID uint,
	status,
	kind string,
) (database.ApplicationStatus, error) {
	if sm.db == nil {
		return database.ApplicationStatus{}, errors.New("database not initialized")
	}

	applicationStatus, err := sm.GetApplicationStatus(nil, statusID)
	if err != nil {
		return database.ApplicationStatus{}, err
	}

	if !sm.IsUserAdmin(userID) && applicationStatus.UserID != userID {
		return database.ApplicationStatus{}, errors.New("you are not allowed to make changes to this application status")
	}

	convKind := database.ApplicationStatusKind(kind)
	if !convKind.Valid() {
		return database.ApplicationStatus{}, errors.New("invalid application kind")
	}

	if status == "" {
		return database.ApplicationStatus{}, errors.New("application status is empty")
	}

	if applicationStatus.Status != status {
		applicationStatus.Status = status
	}

	if applicationStatus.Kind != convKind {
		applicationStatus.Kind = convKind
	}

	result := sm.db.Save(&applicationStatus)

	if result.Error != nil {
		return database.ApplicationStatus{}, result.Error
	} else if result.RowsAffected == 0 {
		return database.ApplicationStatus{}, errors.New("no rows were affected")
	}

	return applicationStatus, nil
}

func (sm *ServiceManager) GetApplicationStatus(tx *gorm.DB, id uint) (database.ApplicationStatus, error) {
	if sm.db == nil {
		return database.ApplicationStatus{}, errors.New("database not initialized")
	}
	db := sm.transactionOrDefault(tx)
	var applicationStatus database.ApplicationStatus
	result := db.
		Where("id = ?", id).
		First(&applicationStatus)

	if result.Error != nil {
		return database.ApplicationStatus{}, result.Error
	}

	return applicationStatus, nil
}

func (sm *ServiceManager) FindApplicationStatusByName(
	userID uuid.UUID,
	name string,
	showArchived bool,
) (database.ApplicationStatus, error) {
	if sm.db == nil {
		return database.ApplicationStatus{}, errors.New("database not initialized")
	}
	var applicationStatus database.ApplicationStatus
	query := sm.db.Where("lower(status) = lower(?) and user_id = ?", name, userID)
	if !showArchived {
		query = query.Where("archived = ?", false)
	}
	result := query.
		First(&applicationStatus)
	if result.Error != nil {
		return database.ApplicationStatus{}, result.Error
	}

	return applicationStatus, nil
}

func (sm *ServiceManager) GetAllApplicationStatus(
	userID uuid.UUID,
	showArchived bool,
) ([]database.ApplicationStatus, error) {
	if sm.db == nil {
		return []database.ApplicationStatus{}, errors.New("database not initialized")
	}
	var applicationStatus []database.ApplicationStatus
	query := sm.db.
		Where("user_id = ?", userID)
	if !showArchived {
		query = query.Where("archived = ?", false)
	}

	result := query.Find(&applicationStatus)

	if result.Error != nil {
		return []database.ApplicationStatus{}, result.Error
	}

	return applicationStatus, nil
}

func (sm *ServiceManager) DeleteApplicationStatus(
	userID uuid.UUID,
	statusID uint,
	softDelete bool,
) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	if sm.IsStatusDefault(statusID, userID) {
		return errors.New("cannot delete a default application status")
	}
	if softDelete {
		jobApplications, err := sm.FindJobApplicationByStatus(userID, statusID)
		if err != nil {
			return err
		}
		err = sm.db.Transaction(func(tx *gorm.DB) error {
			for _, jobApplication := range jobApplications {
				err := sm.DeleteJobApplication(tx, jobApplication.ID, userID)
				if err != nil {
					return err
				}
			}
			var applicationStatus database.ApplicationStatus
			return tx.Delete(&applicationStatus, "id = ? and user_id = ?", statusID, userID).Error
		})
		return err
	}
	status, err := sm.GetApplicationStatus(nil, statusID)
	if err != nil {
		return err
	}
	if !sm.IsUserAdmin(userID) && status.UserID != userID {
		return errors.New("you are not allowed to make changes to this application status")
	}
	status.Archived = true
	result := sm.db.Save(&status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sm *ServiceManager) GetSuggestionForApplicationStatus(
	userID uuid.UUID,
	showArchived bool,
) ([]utils.SuggestionRecord, error) {
	if sm.db == nil {
		return []utils.SuggestionRecord{}, errors.New("database not initialized")
	}
	query := sm.db.Where("user_id = ?", userID)
	if !showArchived {
		query = query.Where("archived = ?", false)
	}
	var data []database.ApplicationStatus
	result := query.Find(&data)
	if result.Error != nil {
		return []utils.SuggestionRecord{}, result.Error
	}
	suggestionRecord := make([]utils.SuggestionRecord, 0, len(data))
	for _, d := range data {
		stringID := strconv.Itoa(int(d.ID))
		temp := utils.SuggestionRecord{

			Label: d.Status,
			Value: stringID,
			// Might be this still not sure so keep it here for a reminder
			//Value: utils.ToSnakeCase(d.Status),
		}
		suggestionRecord = append(suggestionRecord, temp)
	}
	return suggestionRecord, result.Error
}

func (sm *ServiceManager) GetAllApplicationStatusAdmin(showArchived bool) ([]database.ApplicationStatus, error) {
	if sm.db == nil {
		return []database.ApplicationStatus{}, errors.New("database not initialized")
	}
	var applicationStatus []database.ApplicationStatus
	query := sm.db.Where("archived = ?", false)
	if showArchived {
		query = query.Where("archived = ?", true)
	}
	data := query.Find(&applicationStatus)
	if data.Error != nil {
		return []database.ApplicationStatus{}, data.Error
	}
	return applicationStatus, nil
}

func (sm *ServiceManager) DoesUserOwnApplicationStatus(
	userID uuid.UUID,
	statusID uint,
) bool {
	if sm.db == nil {
		return false
	}
	status := database.ApplicationStatus{}
	if sm.IsUserAdmin(userID) {
		return true
	}
	result := sm.db.First(&status, "id = ? and user_id = ?", statusID, userID)
	if result.Error != nil {
		return false
	}
	return true
}

func (sm *ServiceManager) IsStatusArchived(statusID uint) bool {
	if sm.db == nil {
		return false
	}
	status, err := sm.GetApplicationStatus(nil, statusID)
	if err != nil {
		return false
	}
	return status.Archived
}

func (sm *ServiceManager) IsStatusDefault(statusID uint, userID uuid.UUID) bool {
	if sm.db == nil {
		return false
	}
	user, err := sm.GetUserByID(userID)
	if err != nil {
		return false
	}
	if user.UserSettings.DefaultApplicationStatusID == nil {
		return false
	}
	return *user.UserSettings.DefaultApplicationStatusID == statusID
}

func (sm *ServiceManager) UnarchiveJobApplicationStatus(statusID uint, userID uuid.UUID) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	if !sm.DoesUserOwnApplicationStatus(userID, statusID) && !sm.IsUserAdmin(userID) {
		return errors.New("user does not own job application status")
	}

	db := sm.transactionOrDefault(sm.db)
	err := db.Transaction(func(tx *gorm.DB) error {
		status, err := sm.GetApplicationStatus(tx, statusID)
		if err != nil {
			return err
		}
		if !status.Archived {
			return errors.New("job application status is not archived")
		}
		status.Archived = false
		if err := tx.Save(status).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (sm *ServiceManager) CountApplicationStatusByKind(tx *gorm.DB, userID uuid.UUID, byJobs bool) (utils.ApplicationStatusKindCounts, error) {
	if sm.db == nil {
		return utils.ApplicationStatusKindCounts{}, errors.New("database not initialized")
	}
	var rows []utils.KindCount
	var result *gorm.DB
	if !byJobs {
		result = tx.Model(&database.ApplicationStatus{}).
			Select("kind, count(*) as count").
			Where("user_id = ? and archived = false and deleted_at is null", userID).
			Group("kind").
			Find(&rows)
	} else {
		result = tx.Model(&database.ApplicationStatus{}).
			Select("application_statuses.kind, count(*) as count").
			Joins("inner join job_applications on job_applications.status_id = application_statuses.id").
			Where("application_statuses.user_id = ? and job_applications.deleted_at is null", userID).
			Group("application_statuses.kind").
			Find(&rows)
	}

	if result.Error != nil {
		return utils.ApplicationStatusKindCounts{}, result.Error
	}
	byKind := make(map[database.ApplicationStatusKind]int64)
	for _, row := range rows {
		byKind[row.Kind] = row.Count
	}
	counts := utils.ApplicationStatusKindCounts{}
	v := reflect.ValueOf(&counts).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		key, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		v.Field(i).SetInt(byKind[database.ApplicationStatusKind(key)])
	}
	return counts, nil
}
