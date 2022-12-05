# FSTSP-Instances

## Instances
The instances are from TSPLIB repository. All edges are the type EUC_2D.

## Initial Solution
The initial solution is a TSP tour.
TSP solution was obtained by changing edge type from EUC_2D to MAN_2D, then running Concorde to obtain the optimal solution considering Manhattan distance.

## Tau File
The tau.csv file defines the truck travel time from one node to another, considering a speed of 40km/h.

## Tauprime File
The tauprime.csv file defines the drone travel time from one node to another, considering a speed of 40km/h.

## Node File
The nodes.csv file defines the coordinates of each customer, in the format:
                        Node#  x    y    isTooHeavy
                          0                    0
                          1                  {0,1}
                         c+1                   0

Node is the id of the customer.
x -> x-coordinate of the node.
y -> y-coordinate of the node.
isTooHeavy -> 1 The customer's parcel is too heavy for the drone.
              0 The drone can serve this customer.


