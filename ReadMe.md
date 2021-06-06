#Data Models
1. Guest
    { id, guestname, accompanying_guests, tableId }
   What : To represent a guest entity.
2. GuestLog
   { id, guestid, ispresent, time_arrived, time_departed, accompanying_guest }
   What : Entity representing the guests arrived at and left from party.
3. Venue
   { id, venuename }
   What : Entity to provide venue details.
4. TableInfo
   { id, state, capacity, venueid }
   What : Entity for a table that should be mapped to a venue.
   
##Workflows
1. Add Guest to guestlist.
   API : `POST /guest_list/name`. Create a record in Guest Table.
2. Fetch the guestlist.
   API : `GET /guest_list`. Read all records from Guest table and return.
3. Allow/Disallow a guest on arrival.
   API : `PUT /guests/name`. Add/Update the record in GuestLog table.  
4. Mark departure of guest.
   API : `DELETE /guests/name`. Mark departure by updating the ispresent field of the record.
5. Fetch the guests present in the party (along with few other details).
   API : `GET /guests`. Get the records from GuestLog that have ispresent field set to 1 i.e. they are present in the venue and return.
6. Get the empty seat count.
   API : `GET /seats_empty`. Get the totalOccupancy by summing up the table's capacity from TableInfo entity, Get the occupied seats from GuestLog entity based on ispresent field.
   Calculate the difference and return.
   
###File Significance
1. _routes.go_ : To map the APIs to respective methods.
2. _Controller.go_ : Common interface that will route the incoming request to respective Controller
3. _GuestController.go_ : This controller will interact with Guest and GuestLog Entity representative packages.
4. _TableController_ : This controller will the interface to communicate to TableInfo entity.
5. _Guest.go_ : This is Entity model of Guest Table. It wraps all the operations done on Guest Entity.
6. _GuestLog.go_ : This is Entity model of GuestLog Table. It wraps all the operations done on GuestLog Entity.
7. _Table.go_ : This is Entity model of TableInfo table. It wraps all the operations done on TableInfo Entity.
8. _server.go_ : This the **entry point of the service** containing main() method.

###DB Setup
1. _partydbConsoleSetup.sql_ : This contains sql scripts to setup the DB for the project to work.


##Areas Covered
1. Basic implementation of party management. 
    - For adding guests, fetching guest details
    - For tracking guests entry, exit
    
##Not Covered 
1. To keep the project size small and complexity at minimal level, not covered following :
    - Any Concurrency and go routines.
      
    - No implementation to add/remove TableInfo entries. That can be added using mysql server separately
    - Any standard patterns to keep code extensible.
    - Users scenario to handle admin, guests, venue manager etc.
    
