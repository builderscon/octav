//go:generate gendb -t Conference -t Room -t Session -t User -t Venue -t LocalizedString -d db
//go:generate genmodel -t Session -t Conference -t Room -t User -t Venue -d .
//go:generate gentransport -t CreateSessionRequest -t UpdateSessionRequest -t ListVenueRequest -t ListSessionsByConferenceRequest -t UpdateConferenceRequest -d .

package octav
