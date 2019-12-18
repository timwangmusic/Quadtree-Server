# Quadtree-Server

## Rest API
* Add place
   * path: `/v1/addplace`
   * A POST method for inserting POI into the database

* Search Places
   * path: `/v1/searchplaces`
   * A GET method for performing POI range search

## In-memory
* Persistence is not implemented yet. Restart server flushes all the POI in the server.
