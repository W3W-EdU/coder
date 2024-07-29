package codersdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

type NotificationsSettings struct {
	NotifierPaused bool `json:"notifier_paused"`
}

type NotificationTemplate struct {
	ID            uuid.UUID `json:"id" format:"uuid"`
	Name          string    `json:"name"`
	TitleTemplate string    `json:"title_template"`
	BodyTemplate  string    `json:"body_template"`
	Actions       string    `json:"actions" format:""`
	Group         string    `json:"group"`
	Method        string    `json:"method"`
	Kind          string    `json:"kind"`
}

type NotificationPreference struct {
	NotificationTemplateID uuid.UUID `json:"id" format:"uuid"`
	Disabled               bool      `json:"disabled"`
	UpdatedAt              time.Time `json:"updated_at" format:"date-time"`
}

// GetNotificationsSettings retrieves the notifications settings, which currently just describes whether all
// notifications are paused from sending.
func (c *Client) GetNotificationsSettings(ctx context.Context) (NotificationsSettings, error) {
	res, err := c.Request(ctx, http.MethodGet, "/api/v2/notifications/settings", nil)
	if err != nil {
		return NotificationsSettings{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return NotificationsSettings{}, ReadBodyAsError(res)
	}
	var settings NotificationsSettings
	return settings, json.NewDecoder(res.Body).Decode(&settings)
}

// PutNotificationsSettings modifies the notifications settings, which currently just controls whether all
// notifications are paused from sending.
func (c *Client) PutNotificationsSettings(ctx context.Context, settings NotificationsSettings) error {
	res, err := c.Request(ctx, http.MethodPut, "/api/v2/notifications/settings", settings)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotModified {
		return nil
	}
	if res.StatusCode != http.StatusOK {
		return ReadBodyAsError(res)
	}
	return nil
}

// UpdateNotificationTemplateMethod modifies a notification template to use a specific notification method, overriding
// the method set in the deployment configuration.
func (c *Client) UpdateNotificationTemplateMethod(ctx context.Context, notificationTemplateId uuid.UUID, method string) error {
	res, err := c.Request(ctx, http.MethodPut,
		fmt.Sprintf("/api/v2/notifications/templates/%s/method", notificationTemplateId),
		UpdateNotificationTemplateMethod{Method: method},
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotModified {
		return nil
	}
	if res.StatusCode != http.StatusOK {
		return ReadBodyAsError(res)
	}
	return nil
}

// GetSystemNotificationTemplates retrieves all notification templates pertaining to internal system events.
func (c *Client) GetSystemNotificationTemplates(ctx context.Context) ([]NotificationTemplate, error) {
	res, err := c.Request(ctx, http.MethodGet, "/api/v2/notifications/templates/system", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, ReadBodyAsError(res)
	}

	var templates []NotificationTemplate
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, xerrors.Errorf("read response body: %w", err)
	}

	if err := json.Unmarshal(body, &templates); err != nil {
		return nil, xerrors.Errorf("unmarshal response body: %w", err)
	}

	return templates, nil
}

// GetUserNotificationPreferences TODO
func (c *Client) GetUserNotificationPreferences(ctx context.Context) ([]NotificationPreference, error) {
	res, err := c.Request(ctx, http.MethodGet, "/api/v2/notifications/preferences", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, ReadBodyAsError(res)
	}

	var prefs []NotificationPreference
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, xerrors.Errorf("read response body: %w", err)
	}

	if err := json.Unmarshal(body, &prefs); err != nil {
		return nil, xerrors.Errorf("unmarshal response body: %w", err)
	}

	return prefs, nil
}

// UpdateUserNotificationPreferences TODO
func (c *Client) UpdateUserNotificationPreferences(ctx context.Context, req UpdateUserNotificationPreferences) ([]NotificationPreference, error) {
	res, err := c.Request(ctx, http.MethodPut, "/api/v2/notifications/preferences", req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, ReadBodyAsError(res)
	}

	var prefs []NotificationPreference
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, xerrors.Errorf("read response body: %w", err)
	}

	if err := json.Unmarshal(body, &prefs); err != nil {
		return nil, xerrors.Errorf("unmarshal response body: %w", err)
	}

	return prefs, nil
}

type UpdateNotificationTemplateMethod struct {
	Method string `json:"method,omitempty" example:"webhook"`
}

type UpdateUserNotificationPreferences struct {
	TemplateDisabledMap map[string]bool `json:"template_disabled_map"`
}
