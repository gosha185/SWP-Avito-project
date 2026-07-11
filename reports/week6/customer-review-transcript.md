# Customer Review Transcript – Week 5

**Date:** 11 July 2026  
**Participants:** Industry Expert (Customer), Team Members    
**Meeting purpose:** Sprint Review – demonstration of Sprint 4 results, discussing trial verse and product delivery, 
and feedback on the team's development process.

---

**Team Member:**
Let's start. We are currently in the fourth sprint. This week our goal was to discuss transitioning the product to you – the handover. Next week the actual full transfer will take place, and the product will be on your side. Let's begin with a screen share. Stepan prepared the customer handover documentation, he’ll walk you through it shortly. We also added a link to the handover in the Readme and in Contributing Agents MD. Our sprint backlog for this week looks like this: most tasks are closed, some will be closed after today’s meeting. We added DORs and DODs to issue templates as you suggested, and fixed the error – we actually fixed it last week, but now it’s final. Let’s start by discussing the product transition. We have a few questions for you. First, is the product ready for handover? I know our backend team is still thinking about optimisation, but tests are already written. So overall the product is ready. We also want to ask: is the customer already using the product? Obviously not, since we haven’t transferred it yet. And we would like to know what methods you prefer for the transfer.

**Expert:**
We need to transfer not only the product but also documentation and architecture. So we need to gather ADRs, information about implemented user stories, and confirmation from QA that everything is tested and working. Also we’ll need a formal Docker container. The project is somewhat educational, so we’ll take it, but whether it will be actually deployed is uncertain. The question is: is formal deployment critical for you to pass the course? Is that a hard requirement? If so, we’ll figure something out.

**Team Member:**
We need to make sure you have no problems with your diploma, so that you can integrate the product into your system.

**Expert:**
Okay, formally you have a Docker container – that is a reasonably repeatable way to deliver a working service. So if you have a Docker container, that’s a sensible default.

**Team Member:**
Alright, thank you. So next week we’ll have the final handover. We’ll be available if you run into any issues with the transfer. Now we can show the product again.

**Team Member:**
Yes, let’s do that.

**Team Member:**
Stepan, can you show the fixed error right now?

**Team Member:**
Yes, I’m sharing my screen. Can you see it? Last time the bug didn’t appear fixed because we pushed but didn't deploy – that’s why the server still returned 500. Let me quickly go through the steps. So, a batch is created. Hold. Now, I tried to hold this order ID, but I already used it for a hold in the previous demo video. So I need to change it to, say, 8, and then it returns a successful creation. I wanted to ask: should we fix the case where trying to hold the same order again returns 500? Should it return a different status?

**Team Member:**
Well, we are considering that this order has already been held before. But the idempotency key is different, right? So the previous hold used a different key. A 500 error is never good – it means we didn’t handle this scenario. Logically, we should treat it as a 400 error.

**Expert:**
So it's a bad request, because the service rules you designed should not allow double point assignments on the same order. So I would change it from 500 to 400. Beyond just a stylistic improvement, there's a more important reason: services under heavy load are usually monitored. 500s are seen as unexpected behaviour. 400s are expected errors – they are okay, they are predicted scenarios. But 500s trigger alerts, wake up on-call engineers, and ruin their days off. So we try to avoid 500s. I’d definitely change it to 400.

**Team Member:**
Okay, thanks. And for confirmation – there was a bug where confirming twice returned 500, right? Now, first confirmation works, second time returns bad request – hold not found. I also fixed the same issue for cancel.

**Team Member:**
Okay, good.

**Expert:**
Again, I’m not asking you to redo everything, just sharing some thoughts. For the second confirmation of a hold, you could process it differently – if you receive the same idempotency key for a second confirmation, you could return 200, because the hold was already confirmed once. A 400 could cause problems if the client retries aggressively: the first request succeeded but the connection dropped, the retry sends the same key, and you return 400 – then the client might think the request was malformed, while it just wants to be sure the hold went through. Since you already specify 409 in your spec, you could treat 409 as success in retry logic, but it would be simpler to return 200 on retry. I’m not insisting on changes, just keep in mind that it's worth thinking about consumer convenience.

**Team Member:**
Thank you. Let me show you the customer handover document quickly. It describes the product status, functionality, FIFO, idempotency. Should we go through it in detail or just skim?

**Expert:**
No, let's skim. I see you have...

**Team Member:**
Access to the current product – the deployed one. Instructions on how to run locally, via environment variables, as you mentioned. Here are all required env values, with defaults, and a warning not to commit them. And here is deployment on a server, and how to run migrations.

**Expert:**
Good, that’s correct.

**Team Member:**
And here are limitations – we planned to add the limits from tests, but it's still in progress. The situations are described.

**Team Member:**
It looks good.

**Expert:**
Let's see what’s in the documentation. Architecture is there, migrations, okay. Development process, okay. Testing… Also, it would be good to have user stories – you have a project and architecture, but user stories describe what it does end-to-end. And you probably have test status for each story, right?

**Team Member:**
I’m not sure about that.

**Expert:**
Actually, I think I’ve seen them. You should attach that too. As a customer, I want to know what was actually implemented and what works, confirmed by QA. So please include that.

**Team Member:**
Okay. I also created CONTRIBUTING.md and AGENTS.md – I asked an AI to generate the agents instructions. CONTRIBUTING is basically a short summary. It also mentions testing and Git workflow. Here’s a template for pull requests.

**Team Member:**
Is the PR template automatically applied?

**Team Member:**
Yes, it is.

**Expert:**
It doesn’t hurt to have it, even though it’s already configured. Okay, good.

**Team Member:**
(Unclear question about "atom"?)

**Expert:**
It seems you have almost everything ready. By the way, do we have one or two formal meetings left? Two, right?

**Team Member:**
Yes, yes.

**Expert:**
Is it the 22nd? Actually it’s the 21st? Sorry. So next week we formally lock the handover. Let’s meet on Thursday – Friday afternoon is also possible but not too late because I’ll be travelling. Saturday I’ll be travelling too. So Thursday works?

**Team Member:**
Yes, Thursday is fine.

**Expert:**
Great, we’ll agree on the exact time soon.

**Team Member:**
Quick recap: next week we have the sprint preview, we complete the handover, prepare for the final presentation, and we only have to fix the issues you pointed out today.

**Team Member:**
Okay, yes.

**Team Member:**
Thank you for your time. One more question about the delivery format: we considered either a link to a public repository or a zip archive.

**Expert:**
A repository link is more convenient, but if you prefer an archive, that’s fine too.

**Team Member:**
We were thinking about confidentiality – who would have access to the repo. If it’s public, then just a link. Also, our repository currently contains weekly reports and course assignments – should we include those or clean them up? I’ll think about it. Is it currently public or private?

**Team Member:**
It’s public.

**Expert:**
Let’s leave it as is. I don’t think having course assignments is a downside – it’s normal. More information is better than less, unless there are university requirements to remove them. If not, keep it.

**Team Member:**
Alright, we’ll keep it. Thanks. So, have we covered everything?

**Team Member:**
Yes, I think so. Thank you very much, goodbye.

**Expert:**
Just so you know, generating the link sometimes takes an hour or two, so I’ll try to send it to you this evening. If I forget, ping me and I’ll check and send it. Have a great day.

**Team Member:**
Good luck. Have a nice day, thank you. Goodbye.
