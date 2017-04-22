package services

//
// ArangoDB Documentation
// https://docs.arangodb.com/3.1/cookbook/AQL/MigratingEdgeFunctionsTo3.html
// Using as anonymous graph
//  FOR v, e IN 1..1 OUTBOUND @startId @@edgeCollection OPTIONS {bfs: true, uniqueVertices: 'global'}
//  RETURN {edge: e, vertex: v}
//

const (
	// Edge collections
	PRJ_EDGES = "prjedges"

	// Named graphs not used yet. Used anonymous graphs
	//PRJ_GRAPH = "prjusers"

	// Roles
	ROLE_PO = "PO" 			// Product Owner
	ROLE_RTE = "RTE" 		// Release Train Engineer
	ROLE_ARCHITECTOR = "ARCHITECT" 	// Architect
	ROLE_BA = "BA"			// Business Analyst
	ROLE_SA = "SA"			// System Analyst
	ROLE_SM = "SM"			// Scram Master
	ROLE_DEV = "DEV"		// Developer
	ROLE_BP = "BP"			// Business Partner
	ROLE_TPM = "TPM"		// Technical Project Manager
	ROLE_PM = "PM"			// Project Manager
	ROLE_VSE = "VSE"		// Something else

	// Statuses
	PROCESS_STATUS = "("
)
