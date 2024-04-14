/*
   DIEGO - A data importer extension for Hugo
   Copyright (C) 2024 Vinícius Moraes <vinicius.moraes@eternodevir.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package youtube

import "fmt"

// Subscriptions represents the official input data fields.
type Subscriptions struct {
	Name       string
	ChannelID  string
	ChannelURL string
}

// BindFile binds the CSV values into the Subscriptions struct.
func (y *Subscriptions) BindFile(record *[]string) error {
	y.Name = (*record)[2]
	y.ChannelID = (*record)[0]
	y.ChannelURL = (*record)[1]

	return nil
}

// FetchFromHTTP gets additional values from the URLs for the Subscriptions struct.
func (y *Subscriptions) FetchFromHTTP() error {
	return nil
}

// BindHTML generates the Hugo shortcode for the Subscriptions struct.
func (y *Subscriptions) BindHTML(shortcode, comment *string, model string) error {
	return htmlSubscriptions(shortcode, comment, model)
}

// SubscriptionsComplete represents the official input data fields.
type SubscriptionsComplete struct {
	Name       string
	ChannelID  string
	ChannelURL string
}

// BindFile binds the CSV values into the SubscriptionsComplete struct.
func (y *SubscriptionsComplete) BindFile(record *[]string) error {
	y.Name = (*record)[2]
	y.ChannelID = (*record)[0]
	y.ChannelURL = (*record)[1]

	return nil
}

// FetchFromHTTP gets additional values from the URLs for the SubscriptionsComplete struct.
func (y *SubscriptionsComplete) FetchFromHTTP() error {
	return nil
}

// BindHTML generates the Hugo shortcode for the SubscriptionsComplete struct.
func (y *SubscriptionsComplete) BindHTML(shortcode, comment *string, model string) error {
	return htmlSubscriptions(shortcode, comment, model)
}

// Common to both structs.
func htmlSubscriptions(shortcode, comment *string, model string) error {
	htmlTemplate := `%s
<table>
  <tbody>
    {{ range sort .Site.Data.%s "name" }}
    <tr>
      <td>
	<strong>{{ .channelid }}</strong>
      </td>
      <td>
	<a href="{{ .channelurl }}">YouTube</a>
      </td>
    </tr>
    {{ end }}
  </tbody>
</table>`

	*shortcode = fmt.Sprintf(htmlTemplate, *comment, model)

	return nil
}
