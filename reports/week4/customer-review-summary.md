# Customer Review Summary – Week 4

**Date:** 25 June 2026

## Participants

* **Customer:** Expert from Avito
* **Team:** Ekaterina, Ivan, Georgii, and other team members

## Sprint Goal Review

The Sprint Goal was to complete the missing functionality from the previous sprint, stabilize the product, and ensure
that it successfully passed automated tests. The team confirmed that the planned work had been completed and that the
product was now operational.

The customer also explained the purpose of a Sprint Goal, emphasizing that it should describe the most important
business objective of the sprint rather than simply repeat the Sprint Backlog. A clear Sprint Goal helps the team
prioritize work when unexpected issues arise.

## Delivered Increment

The team presented the current MVP through Swagger because a graphical user interface has not yet been implemented.

The demonstrated functionality included:

* awarding bonus points to users;
* retrieving user balances;
* querying points by expiration period;
* placing points on hold;
* confirming and cancelling holds;
* viewing transaction history.

The customer noted that the demonstrated functionality represented meaningful progress toward a usable product
increment.

## User Acceptance Testing (UAT)

The customer observed several user scenarios performed during the demonstration:

* awarding bonus points with expiration dates;
* checking balance calculations;
* creating and confirming holds;
* cancelling holds and restoring available balance;
* viewing transaction history;
* reserving points originating from multiple award transactions.

The demonstrated scenarios behaved as expected overall.

One issue was identified: attempting to cancel a non-existent hold returned **HTTP 500** instead of the expected **HTTP
400 Bad Request**.

## Quality Evidence

The team reported that:

* most unit tests had been implemented;
* concurrent request tests had been added;
* the product successfully passed the available automated tests;
* Sprint progress was tracked using a board with **To Do**, **In Progress**, **Review**, and **Done** columns together
  with Daily Scrums.

## Customer Feedback

The customer provided several recommendations:

* explain Sprint Reviews from a business perspective rather than focusing only on technical implementation;
* present the Sprint Goal, achieved results, demonstration, and next Sprint Goal during each Sprint Review;
* aim to deliver complete, usable increments instead of partially completed features;
* avoid assuming that current business requirements will remain unchanged;
* return more appropriate HTTP status codes for invalid client requests;
* define Sprint Goals before each sprint begins and use them to guide prioritization;
* consider returning either an empty response or echoed request data for successful mutation endpoints instead of
  internal implementation details.

## Requested Changes

* Return **HTTP 400 Bad Request** instead of **HTTP 500 Internal Server Error** when cancelling a non-existent hold.

## Decisions

* The demonstrated functionality was accepted as a valid Sprint increment.
* The current API design is acceptable for the existing requirements.
* Swagger will continue to be used for demonstrations until a UI becomes available.

## Action Points

* Fix incorrect HTTP status codes for invalid operations.
* Add additional UAT scenarios for expiration priority and multi-transaction reservations.
* Define the Sprint Goal before Sprint Planning.
* Continue improving the product according to customer feedback and upcoming university requirements.
* Verify that points with the earliest expiration date are consumed first.
* Test scenarios involving multiple point awards and reservations across different transactions.

