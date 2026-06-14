Customer Meeting Summary 

Date: 11 June 2026  
Attendees:  
Customer representative: Representative of Avito  
Interviewers:  
   Stepan Grechinskii  
   Georgii Sergeev (asked most of the questions)  
   Leilia Zakirova  
   Ekaterina Efremova  
   Ivan Shpenkov  

Artifacts presented:  
 Draft user stories  
 Proposed MVP v1 scope  
 OpenAPI specification and Swagger UI (on VM `10.93.26.189:8080`)  
 Postman collection  

Discussion points:

1. API requirements. The customer clarified that the system must provide:
    Current available balance.
    Held points per concrete order.
    Points expiring within a configurable window – only aggregate, not per‑day breakdown.

2. Idempotency and failure handling. The customer strongly recommended using idempotency keys to prevent duplicate accruals when retries occur. The vague requirement “bonuses not lost during failures” should be replaced with a testable one: “no duplicate bonus accruals on retries”.

3. Hold management. Only full confirmation or full release of a hold is required. Automatic release after a configurable TTL (e.g., 24 hours) is acceptable; the TTL could be made per‑request.

4. User story priorities:
    Showing only available balance to end‑user → Must Have.
    Showing points expiring soon → Could Have for MVP.
    Admin story “total points in circulation” → Should Have.
    The refined idempotency story → Should Have.

5. Expiration and spending order. The customer suggested FEFO. 

6. Automatic hold release TTL. The customer agreed that a configurable TTL (e.g., 1 hour for banking transactions, 24 hours for orders) is a good enhancement.

Decisions and approvals:

 The customer approved the user stories after adjustments.  
 The MoSCoW priorities were accepted with the changes (US‑05 moved to Could Have, US‑07 added as Must Have, US‑02 rewritten and moved to Should Have).  
 The initial MVP v1 scope was approved as:    
   US‑06 (view points held for a concrete order)  
   US‑07 (view current available point balance)  
 The customer gave written consent to the MIT‑licensed public development model before repository creation.  
 Permission to publish a sanitised meeting transcript in the repository was granted.
