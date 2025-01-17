package athenahealth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Appointment struct {
	AppointmentID              string `json:"appointmentid"`
	AppointmentStatus          string `json:"appointmentstatus"`
	AppointmentType            string `json:"appointmenttype"`
	AppointmentTypeID          string `json:"appointmenttypeid"`
	ChargeEntryNotRequired     bool   `json:"chargeentrynotrequired"`
	Date                       string `json:"date"`
	DepartmentID               string `json:"departmentid"`
	Duration                   int    `json:"duration"`
	EncounterID                string `json:"encounterid"`
	PatientAppointmentTypeName string `json:"patientappointmenttypename"`
	ProviderID                 string `json:"providerid"`
	StartTime                  string `json:"starttime"`
}

// GetAppointment - Single appointment.
// GET /v1/{practiceid}/appointments/{appointmentid}
// https://developer.athenahealth.com/docs/read/appointments/Appointments#section-1
func (h *HTTPClient) GetAppointment(ctx context.Context, id string) (*Appointment, error) {
	out := []*Appointment{}

	_, err := h.Get(ctx, fmt.Sprintf("/appointments/%s", id), nil, &out)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("Unexpected length returned")
	}

	return out[0], nil
}

type AppointmentCustomField struct {
	CaseSensitive  bool   `json:"casesensitive"`
	CustomFieldID  int    `json:"customfieldid"`
	DisallowUpdate bool   `json:"disallowupdate"`
	Name           string `json:"name"`
	Searchable     bool   `json:"searchable,omitempty"`
	Select         bool   `json:"select"`
	Type           string `json:"type"`
	SelectList     []struct {
		OptionValue string `json:"optionvalue"`
		OptionID    int    `json:"optionid"`
	} `json:"selectlist,omitempty"`
}

type listAppointmentCustomFieldsResponse struct {
	AppointmentCustomFields []*AppointmentCustomField `json:"appointmentcustomfields"`
}

// ListAppointmentCustomFields - List of appointment custom fields (practice specific).
// GET /v1/{practiceid}/appointments/customfields
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Custom_Fields#section-0
func (h *HTTPClient) ListAppointmentCustomFields(ctx context.Context) ([]*AppointmentCustomField, error) {
	out := &listAppointmentCustomFieldsResponse{}

	_, err := h.Get(ctx, "/appointments/customfields", nil, &out)
	if err != nil {
		return nil, err
	}

	return out.AppointmentCustomFields, nil
}

type BookedAppointment struct {
	AppointmentID    string `json:"appointmentid"`
	AppointmentCopay struct {
		CollectedForOther       int `json:"collectedforother"`
		CollectedForAppointment int `json:"collectedforappointment"`
		InsuranceCopay          float32 `json:"insurancecopay"`
	} `json:"appointmentcopay"`
	AppointmentNotes []struct {
		DisplayOnSchedule bool   `json:"displayonschedule"`
		Text              string `json:"text"`
		ID                int    `json:"id"`
	} `json:"appointmentnotes"`
	AppointmentStatus          string `json:"appointmentstatus"`
	AppointmentType            string `json:"appointmenttype"`
	AppointmentTypeID          string `json:"appointmenttypeid"`
	CancelledBy                string `json:"cancelledby"`
	CancelledDatetime          string `json:"cancelleddatetime"`
	CancelReasonID             string `json:"cancelreasonid"`
	CancelReasonName           string `json:"cancelreasonname"`
	CancelReasonNoShow         bool   `json:"cancelreasonnoshow"`
	CancelReasonSlotAvailable  bool   `json:"cancelreasonslotavailable"`
	ChargeEntryNotRequired     bool   `json:"chargeentrynotrequired"`
	CoordinatorEnterprise      bool   `json:"coordinatorenterprise"`
	Copay                      float32 `json:"copay"`
	Date                       string `json:"date"`
	DepartmentID               string `json:"departmentid"`
	Duration                   int    `json:"duration"`
	EncounterID                string `json:"encounterid"`
	HL7ProviderID              int    `json:"hl7providerid"`
	LastModified               string `json:"lastmodified"`
	LastModifiedBy             string `json:"lastmodifiedby"`
	PatientAppointmentTypeName string `json:"patientappointmenttypename"`
	PatientID                  string `json:"patientid"`
	ProviderID                 string `json:"providerid"`
	ScheduledBy                string `json:"scheduledby"`
	ScheduledDatetime          string `json:"scheduleddatetime"`
	StartTime                  string `json:"starttime"`
	TemplateAppointmentID      string `json:"templateappointmentid"`
	TemplateAppointmentTypeID  string `json:"templateappointmenttypeid"`
}

type ListBookedAppointmentsOptions struct {
	DepartmentID      string
	EndDate           time.Time
	PatientID         string
	ProviderID        string
	StartDate         time.Time
	AppointmentStatus string

	Pagination *PaginationOptions
}

type ListBookedAppointmentsResult struct {
	BookedAppointments []*BookedAppointment
	Pagination         *PaginationResult
}

type listBookedAppointmentsResponse struct {
	Appointments []*BookedAppointment `json:"appointments"`

	PaginationResponse
}

// ListBookedAppointments - Booked appointment slots.
// GET /v1/{practiceid}/appointments/booked
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Slots#section-3
func (h *HTTPClient) ListBookedAppointments(ctx context.Context, opts *ListBookedAppointmentsOptions) (*ListBookedAppointmentsResult, error) {
	out := &listBookedAppointmentsResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.ProviderID) > 0 {
			q.Add("providerid", opts.ProviderID)
		}

		if len(opts.DepartmentID) > 0 {
			q.Add("departmentid", opts.DepartmentID)
		}

		if len(opts.PatientID) > 0 {
			q.Add("patientid", opts.PatientID)
		}

		if !opts.StartDate.IsZero() {
			q.Add("startdate", opts.StartDate.Format("01/02/2006"))
		}

		if !opts.EndDate.IsZero() {
			q.Add("enddate", opts.EndDate.Format("01/02/2006"))
		}

		if len(opts.AppointmentStatus) > 0 {
			q.Add("appointmentstatus", opts.AppointmentStatus)
		}

		if opts.Pagination != nil {
			if opts.Pagination.Limit > 0 {
				q.Add("limit", strconv.Itoa(opts.Pagination.Limit))
			}

			if opts.Pagination.Offset > 0 {
				q.Add("offset", strconv.Itoa(opts.Pagination.Offset))
			}
		}
	}

	_, err := h.Get(ctx, "/appointments/booked", q, out)
	if err != nil {
		return nil, err
	}

	return &ListBookedAppointmentsResult{
		BookedAppointments: out.Appointments,
		Pagination:         makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}

type ListChangedAppointmentsOptions struct {
	DepartmentID               string
	LeaveUnprocessed           bool
	PatientID                  string
	ProviderID                 string
	ShowPatientDetail          bool
	ShowProcessedEndDatetime   time.Time
	ShowProcessedStartDatetime time.Time
}

type listChangedAppointmentsResponse struct {
	ChangedAppointments []*BookedAppointment `json:"appointments"`
}

// ListChangedAppointments - Changed appointment slots.
// GET /v1/{practiceid}/appointments/changed
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Slots#section-5
func (h *HTTPClient) ListChangedAppointments(ctx context.Context, opts *ListChangedAppointmentsOptions) ([]*BookedAppointment, error) {
	out := &listChangedAppointmentsResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.ProviderID) > 0 {
			q.Add("providerid", opts.ProviderID)
		}

		if len(opts.DepartmentID) > 0 {
			q.Add("departmentid", opts.DepartmentID)
		}

		if len(opts.PatientID) > 0 {
			q.Add("patientid", opts.PatientID)
		}

		if opts.ShowPatientDetail {
			q.Add("showpatientdetail", strconv.FormatBool(opts.ShowPatientDetail))
		}

		if !opts.ShowProcessedEndDatetime.IsZero() {
			q.Add("showprocessedenddatetime", opts.ShowProcessedEndDatetime.Format("01/02/2006 15:04:05"))
		}

		if !opts.ShowProcessedStartDatetime.IsZero() {
			q.Add("showprocessedstartdatetime", opts.ShowProcessedStartDatetime.Format("01/02/2006 15:04:05"))
		}

		if opts.LeaveUnprocessed {
			q.Add("leaveunprocessed", strconv.FormatBool(opts.LeaveUnprocessed))
		}
	}

	_, err := h.Get(ctx, "/appointments/changed", q, out)
	if err != nil {
		return nil, err
	}

	return out.ChangedAppointments, nil
}

type CreateAppointmentNoteOptions struct {
	AppointmentID     string
	DisplayOnSchedule bool
	NoteText          string
}

// CreateAppointmentNote - Notes for this appointment.
// POST /v1/{practiceid}/appointments/{appointmentid}/notes
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Notes#section-0
func (h *HTTPClient) CreateAppointmentNote(ctx context.Context, appointmentID string, opts *CreateAppointmentNoteOptions) error {
	var form url.Values

	if opts != nil {
		form = url.Values{}

		if len(opts.AppointmentID) > 0 {
			form.Add("appointmentid", opts.AppointmentID)
		}

		if opts.DisplayOnSchedule {
			form.Add("displayonschedule", strconv.FormatBool(opts.DisplayOnSchedule))
		}

		if len(opts.NoteText) > 0 {
			form.Add("notetext", opts.NoteText)
		}
	}

	_, err := h.PostForm(ctx, fmt.Sprintf("/appointments/%s/notes", appointmentID), form, nil)
	if err != nil {
		return err
	}

	return nil
}

type AppointmentNote struct {
	Created           string `json:"created"`
	CreatedBy         string `json:"createdby"`
	DisplayOnSchedule bool   `json:"displayonschedule"`
	NoteID            string `json:"noteid"`
	NoteText          string `json:"notetext"`
}

type ListAppointmentNotesOptions struct {
	AppointmentID string
	ShowDeleted   bool
}

type listAppointmentNotesResponse struct {
	Notes []*AppointmentNote `json:"notes"`
}

// ListAppointmentNotes - Notes for this appointment.
// GET /v1/{practiceid}/appointments/{appointmentid}/notes
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Notes#section-1
func (h *HTTPClient) ListAppointmentNotes(ctx context.Context, appointmentID string, opts *ListAppointmentNotesOptions) ([]*AppointmentNote, error) {
	out := &listAppointmentNotesResponse{}

	q := url.Values{}

	if opts != nil {
		if len(opts.AppointmentID) > 0 {
			q.Add("appointmentid", opts.AppointmentID)
		}

		if opts.ShowDeleted {
			q.Add("showdeleted", strconv.FormatBool(opts.ShowDeleted))
		}
	}

	_, err := h.Get(ctx, fmt.Sprintf("/appointments/%s/notes", appointmentID), q, out)
	if err != nil {
		return nil, err
	}

	return out.Notes, nil
}

type UpdateAppointmentNoteOptions struct {
	AppointmentID     string
	DisplayOnSchedule bool
	NoteID            string
	NoteText          string
}

// UpdateAppointmentNote - Notes for this appointment.
// PUT /v1/{practiceid}/appointments/{appointmentid}/notes/{noteid}
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Notes#section-3
func (h *HTTPClient) UpdateAppointmentNote(ctx context.Context, appointmentID, noteID string, opts *UpdateAppointmentNoteOptions) error {
	var form url.Values

	if opts != nil {
		form = url.Values{}

		if len(opts.AppointmentID) > 0 {
			form.Add("appointmentid", opts.AppointmentID)
		}

		if opts.DisplayOnSchedule {
			form.Add("displayonschedule", strconv.FormatBool(opts.DisplayOnSchedule))
		}

		if len(opts.NoteID) > 0 {
			form.Add("noteid", opts.NoteID)
		}

		if len(opts.NoteText) > 0 {
			form.Add("notetext", opts.NoteText)
		}
	}

	_, err := h.PutForm(ctx, fmt.Sprintf("/appointments/%s/notes/%s", appointmentID, noteID), form, nil)
	if err != nil {
		return err
	}

	return nil
}

type DeleteAppointmentNoteOptions struct {
	AppointmentID string
	NoteID        string
}

// DeleteAppointmentNote - Notes for this appointment.
// DELETE /v1/{practiceid}/appointments/{appointmentid}/notes/{noteid}
// https://developer.athenahealth.com/docs/read/appointments/Appointment_Notes#section-0
func (h *HTTPClient) DeleteAppointmentNote(ctx context.Context, appointmentID, noteID string, opts *DeleteAppointmentNoteOptions) error {
	var form url.Values

	if opts != nil {
		form = url.Values{}

		if len(opts.AppointmentID) > 0 {
			form.Add("appointmentid", opts.AppointmentID)
		}

		if len(opts.NoteID) > 0 {
			form.Add("noteid", opts.NoteID)
		}
	}

	_, err := h.DeleteForm(ctx, fmt.Sprintf("/appointments/%s/notes/%s", appointmentID, noteID), form, nil)
	if err != nil {
		return err
	}

	return nil
}
