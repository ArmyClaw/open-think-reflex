package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ArmyClaw/open-think-reflex/pkg/models"
	"github.com/rivo/tview"
)

// PatternFormMode represents the mode of the pattern form
type PatternFormMode int

const (
	PatternFormCreate PatternFormMode = iota
	PatternFormEdit
)

// PatternFormPanel provides a form for creating and editing patterns
type PatternFormPanel struct {
	theme       *Theme
	view        *tview.Form
	modal       *tview.Modal
	onSave      func(*models.Pattern)
	onCancel    func()
	mode        PatternFormMode
	existingID  string
}

// NewPatternFormPanel creates a new pattern form panel
func NewPatternFormPanel(theme *Theme, onSave func(*models.Pattern), onCancel func()) *PatternFormPanel {
	p := &PatternFormPanel{
		theme:    theme,
		onSave:   onSave,
		onCancel: onCancel,
		mode:     PatternFormCreate,
	}
	p.buildForm()
	return p
}

// buildForm builds the form UI
func (p *PatternFormPanel) buildForm() {
	// Create the form
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Pattern Form")
	
	var trigger, response, strength, threshold string
	
	form.AddInputField("Trigger:", trigger, 50, nil, func(text string) {
		trigger = text
	})
	
	form.AddInputField("Response:", response, 50, nil, func(text string) {
		response = text
	})
	
	form.AddInputField("Strength (0-100):", strength, 50, nil, func(text string) {
		strength = text
	})
	
	form.AddInputField("Threshold (0-100):", threshold, 50, nil, func(text string) {
		threshold = text
	})
	
	form.AddButton("Save", func() {
		p.savePattern(trigger, response, strength, threshold)
	})
	
	form.AddButton("Cancel", func() {
		if p.onCancel != nil {
			p.onCancel()
		}
	})
	
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetBorderColor(p.theme.Border)
	form.SetBackgroundColor(p.theme.Background)
	
	p.view = form
	p.modal = tview.NewModal().
		SetBackgroundColor(p.theme.Background).
		SetTextColor(p.theme.Text)
}

// SetCreateMode sets the form to create mode
func (p *PatternFormPanel) SetCreateMode() {
	p.mode = PatternFormCreate
	p.existingID = ""
	p.clearFields()
	p.view.SetBorder(true).SetTitle(" Create New Pattern ")
}

// SetEditMode sets the form to edit mode with existing values
func (p *PatternFormPanel) SetEditMode(pattern *models.Pattern) {
	p.mode = PatternFormEdit
	p.existingID = pattern.ID
	p.clearFields()
	
	// Populate fields with existing values
	p.populateFields(pattern)
	p.view.SetBorder(true).SetTitle(fmt.Sprintf(" Edit Pattern: %s ", pattern.Trigger))
}

// clearFields clears all form fields
func (p *PatternFormPanel) clearFields() {
	// Get form items and clear them
	for i := 0; i < p.view.GetFormItemCount(); i++ {
		if input, ok := p.view.GetFormItem(i).(*tview.InputField); ok {
			input.SetText("")
		}
	}
}

// populateFields fills the form with pattern data
func (p *PatternFormPanel) populateFields(pattern *models.Pattern) {
	for i := 0; i < p.view.GetFormItemCount(); i++ {
		item := p.view.GetFormItem(i)
		if input, ok := item.(*tview.InputField); ok {
			switch i {
			case 0:
				input.SetText(pattern.Trigger)
			case 1:
				input.SetText(pattern.Response)
			case 2:
				input.SetText(fmt.Sprintf("%.1f", pattern.Strength))
			case 3:
				input.SetText(fmt.Sprintf("%.1f", pattern.Threshold))
			}
		}
	}
}

// savePattern validates and saves the pattern
func (p *PatternFormPanel) savePattern(trigger, response, strength, threshold string) {
	// Validate required fields
	trigger = strings.TrimSpace(trigger)
	response = strings.TrimSpace(response)
	
	if trigger == "" {
		showErrorModal(p.view, "Trigger is required")
		return
	}
	
	if response == "" {
		showErrorModal(p.view, "Response is required")
		return
	}
	
	// Parse strength
	strengthVal := 50.0
	if s, err := strconv.ParseFloat(strength, 64); err == nil {
		strengthVal = s
		if strengthVal < 0 {
			strengthVal = 0
		}
		if strengthVal > 100 {
			strengthVal = 100
		}
	}
	
	// Parse threshold
	thresholdVal := 50.0
	if t, err := strconv.ParseFloat(threshold, 64); err == nil {
		thresholdVal = t
		if thresholdVal < 0 {
			thresholdVal = 0
		}
		if thresholdVal > 100 {
			thresholdVal = 100
		}
	}
	
	// Create or update pattern
	pattern := &models.Pattern{
		ID:           p.existingID,
		Trigger:      trigger,
		Response:     response,
		Strength:     strengthVal,
		Threshold:    thresholdVal,
		Project:      "default",
		DecayEnabled: false,
	}
	
	if p.onSave != nil {
		p.onSave(pattern)
	}
}

// showErrorModal shows an error message in a modal
func showErrorModal(form *tview.Form, message string) {
	// For now, we'll just log the error
	// In a full implementation, we'd show a proper modal
	fmt.Println("Error:", message)
}

// GetView returns the form view
func (p *PatternFormPanel) GetView() tview.Primitive {
	return p.view
}

// SetTheme updates the theme
func (p *PatternFormPanel) SetTheme(theme *Theme) {
	p.theme = theme
	p.view.SetBorderColor(theme.Border)
	p.view.SetBackgroundColor(theme.Background)
}

// DeleteConfirmModal shows a confirmation dialog for deletion
type DeleteConfirmModal struct {
	theme     *Theme
	modal     *tview.Modal
	onConfirm func()
	onCancel  func()
}

// NewDeleteConfirmModal creates a new delete confirmation modal
func NewDeleteConfirmModal(theme *Theme, onConfirm, onCancel func()) *DeleteConfirmModal {
	d := &DeleteConfirmModal{
		theme:     theme,
		onConfirm: onConfirm,
		onCancel:  onCancel,
	}
	d.buildModal()
	return d
}

// buildModal builds the confirmation modal
func (d *DeleteConfirmModal) buildModal() {
	d.modal = tview.NewModal().
		SetBackgroundColor(d.theme.Background).
		SetTextColor(d.theme.Text).
		SetButtonBackgroundColor(d.theme.Primary).
		SetButtonTextColor(d.theme.Background)
}

// Show shows the delete confirmation modal
func (d *DeleteConfirmModal) Show(trigger string, onConfirm, onCancel func()) {
	d.onConfirm = onConfirm
	d.onCancel = onCancel
	
	d.modal.SetText(fmt.Sprintf("Are you sure you want to delete pattern:\n\n'%s'?\n\nThis action cannot be undone.", trigger))
	d.modal.AddButtons([]string{"Delete", "Cancel"})
	d.modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Delete" && d.onConfirm != nil {
			d.onConfirm()
		} else if d.onCancel != nil {
			d.onCancel()
		}
	})
}

// GetView returns the modal view
func (d *DeleteConfirmModal) GetView() tview.Primitive {
	return d.modal
}

// SetTheme updates the theme
func (d *DeleteConfirmModal) SetTheme(theme *Theme) {
	d.theme = theme
	d.modal.SetBackgroundColor(theme.Background)
	d.modal.SetTextColor(theme.Text)
}

// PatternManager manages pattern CRUD operations in the TUI
type PatternManager struct {
	storage   interface {
		SavePattern(ctx context.Context, pattern *models.Pattern) error
		UpdatePattern(ctx context.Context, pattern *models.Pattern) error
		DeletePattern(ctx context.Context, id string) error
	}
	formPanel  *PatternFormPanel
	deleteModal *DeleteConfirmModal
	theme      *Theme
	onUpdate   func() // callback to refresh the pattern list
}

// NewPatternManager creates a new pattern manager
func NewPatternManager(storage interface {
	SavePattern(ctx context.Context, pattern *models.Pattern) error
	UpdatePattern(ctx context.Context, pattern *models.Pattern) error
	DeletePattern(ctx context.Context, id string) error
}, theme *Theme, onUpdate func()) *PatternManager {
	pm := &PatternManager{
		storage:  storage,
		theme:    theme,
		onUpdate: onUpdate,
	}
	
	pm.formPanel = NewPatternFormPanel(theme, pm.handleSave, pm.handleCancel)
	pm.deleteModal = NewDeleteConfirmModal(theme, pm.handleDeleteConfirm, pm.handleDeleteCancel)
	
	return pm
}

// handleSave handles pattern save (create or update)
func (pm *PatternManager) handleSave(pattern *models.Pattern) {
	ctx := context.Background()
	
	var err error
	if pm.formPanel.mode == PatternFormCreate {
		err = pm.storage.SavePattern(ctx, pattern)
	} else {
		err = pm.storage.UpdatePattern(ctx, pattern)
	}
	
	if err != nil {
		fmt.Printf("Error saving pattern: %v\n", err)
		return
	}
	
	// Refresh pattern list
	if pm.onUpdate != nil {
		pm.onUpdate()
	}
}

// handleCancel handles form cancellation
func (pm *PatternManager) handleCancel() {
	// This would typically switch back to main view
}

// handleDeleteConfirm handles delete confirmation
func (pm *PatternManager) handleDeleteConfirm() {
	// Actual deletion is handled externally
}

// handleDeleteCancel handles delete cancellation
func (pm *PatternManager) handleDeleteCancel() {
	// This would typically switch back to main view
}

// ShowCreateForm shows the create pattern form
func (pm *PatternManager) ShowCreateForm() {
	pm.formPanel.SetCreateMode()
}

// ShowEditForm shows the edit pattern form
func (pm *PatternManager) ShowEditForm(pattern *models.Pattern) {
	pm.formPanel.SetEditMode(pattern)
}

// ShowDeleteConfirm shows the delete confirmation modal
func (pm *PatternManager) ShowDeleteConfirm(pattern *models.Pattern, onConfirm func()) {
	pm.deleteModal.Show(pattern.Trigger, onConfirm, func() {})
}

// GetFormView returns the form view
func (pm *PatternManager) GetFormView() tview.Primitive {
	return pm.formPanel.GetView()
}

// GetDeleteView returns the delete confirmation view
func (pm *PatternManager) GetDeleteView() tview.Primitive {
	return pm.deleteModal.GetView()
}

// SetTheme updates the theme
func (pm *PatternManager) SetTheme(theme *Theme) {
	pm.theme = theme
	pm.formPanel.SetTheme(theme)
	pm.deleteModal.SetTheme(theme)
}

// DeletePattern deletes a pattern by ID
func (pm *PatternManager) DeletePattern(ctx context.Context, id string) error {
	return pm.storage.DeletePattern(ctx, id)
}
