package api

// Service ticket fields selected on related issues of a finding
const fragmentServiceTicket = `
fragment ServiceTicketDetails on ServiceTicket {
  id
  externalId
  name
  url
  integration {
    name
  }
}
`

// Service ticket linked to a finding through its related issues
type ServiceTicket struct {
	ExternalId  string `json:"externalId"`
	Id          string `json:"id"`
	Integration *struct {
		Name string `json:"name"`
	} `json:"integration"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type serviceTicketSource struct {
	ServiceTickets []ServiceTicket `json:"serviceTickets"`
}

// FindingServiceTickets holds the issues of a finding that can carry service tickets.
// It is only populated when the service ticket columns are requested.
type FindingServiceTickets struct {
	RelatedIssues *struct {
		Nodes []serviceTicketSource `json:"nodes"`
	} `json:"relatedIssues,omitempty"`
	PostureIssues   []serviceTicketSource `json:"postureIssues,omitempty"`
	PostureIssuesV2 []serviceTicketSource `json:"postureIssuesV2,omitempty"`
}

// LinkedServiceTickets returns the service tickets of the related and posture issues, deduplicated by ID
func (f FindingServiceTickets) LinkedServiceTickets() []ServiceTicket {
	sources := []serviceTicketSource{}
	if f.RelatedIssues != nil {
		sources = append(sources, f.RelatedIssues.Nodes...)
	}
	sources = append(sources, f.PostureIssues...)
	sources = append(sources, f.PostureIssuesV2...)

	seen := map[string]bool{}
	tickets := []ServiceTicket{}
	for _, source := range sources {
		for _, ticket := range source.ServiceTickets {
			if seen[ticket.Id] {
				continue
			}
			seen[ticket.Id] = true
			tickets = append(tickets, ticket)
		}
	}
	return tickets
}
