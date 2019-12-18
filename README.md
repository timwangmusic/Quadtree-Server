# Quadtree-Server
[![Build Status](https://travis-ci.com/weihesdlegend/Quadtree-Server.svg?branch=master)](https://travis-ci.com/weihesdlegend/Quadtree-Server)

## Design Philosophy
The goal is to design a service that enables efficiently look up of Point of Interests (POI) in an area around a specific location.

Traditional databases struggle on this task because of the two-dimensional index on POI make range search operation slow.
Even with spatial index, which is supported by many mainstream databases, the time-complexity is still high when all POIs
are stored in one table in SQL or one collection in MongoDB. Storing in different tables result in complicated queries and data structures.

Quadtree data structure is ideal for designing data storage with two dimensional index. Quadtree only stores POI in its
leaf nodes and each range search narrows down the subtrees need to be considered quickly.

When the number of POIs in a leaf node crosses a threshold, it `split` to 4 children nodes and delegates the POIs to the children nodes.
A maximum depth is also used for limiting the number of `split` and quadtree structure complexity.
When the tree reaches the maximum depth, there is no more `split`.

We have not considered the solution to the case that many POIs are `clustered` around a few commercial or culture rich regions.
It may make sense to set a higher `split` limit for those regions.

## Rest API
* Add place
   * path: `/v1/addplace`
   * A POST method for inserting POI into the database

* Search Places
   * path: `/v1/searchplaces`
   * A GET method for performing POI range search

## In-memory Storage
* Persistence is not implemented yet. Restart server flushes all the POI in the server.
