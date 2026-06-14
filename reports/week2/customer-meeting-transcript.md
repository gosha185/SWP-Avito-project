Participants: 
Customer Expert (expert from the customer company). 
Interviewer 1, Interviewer 2, Interviewer 3, Interviewer 4, Interviewer 5 - interviewers

Customer Expert: Tell me, do we need to record the meeting, or will you record it yourself, or do you need me to?

Interviewer 2: Well, it's better to record it, but it seems that someone is already going to turn it on.

Customer Expert: I'll put it here too, just in case. Tell me which format is best for today's meeting. Will you tell me what you did and ask questions?

Interviewer 2: Well, we have more questions than answers. Yes, we started to create, and we had a question. The first and most important question is: what exactly will we need to fulfill? What kind of requests can come to us? For example, we need to get some burn points. How should that work? Should we give all the data about all points when they burn, or only about specific ones? The first burn points, the second, the third? You don't always need only the first ones.

Customer Expert: Look, the first thing you need to be able to return is points for a specific user, right? That is, you take a user and you should get information about how many points they have in total, how many of them are held. This can be returned in slightly different ways, but conceptually we want to understand how many are available for use, how many are held, and how many will expire in the near future. For the near future, we are probably interested in returning this via a configurable parameter. For example, someone from the outside tells you: tell me how many points a user will burn in 30 days. Because most likely the UI, which you are not building, but which potentially exists, will show the user: dear user, in the next 30 days you will lose so many points, please spend them. Do we need to return how many points will burn on each specific day? No, you don't need that. We are only interested in an aggregated value: that in the next 30 days, say, 200 points will burn. That is the only interesting thing in this part. If we are interested in getting the transaction history – I don't remember if that is in the task, correct me if I'm wrong.

Interviewer 2: But we need to preserve it from a legal standpoint.

Customer Expert: It seems that you need to implement the absence of double spendings, both for legal reasons and for that purpose. It would be good to store them. But I don't remember if there is a requirement to make them retrievable via API. If not, then if there is time left at the end, you could think about it as the next MVP. But you can probably skip it in the database for now.

Interviewer 2: Well, we would like to have a complete list. If we want to get points as of a specific date, when points expire, what is currently available, what is held at the moment, to accrue points with an expiration time, and to hold points.

Customer Expert: Holding doesn't make sense if we don't have an endpoint to confirm the hold and to cancel the hold. The basic tasks will solve that.

Interviewer 2: Another question: if the server cannot perform an operation due to some technical reasons, do we need to worry about executing it as soon as it becomes available, or does that responsibility lie with the API user?

Customer Expert: That responsibility lies with the API user, but it should be possible. Look, your service may crash either before starting to process the client's request or after. For example, you send a command to accrue 200 points. You sent the command, and the service crashed. In what state did it crash? Did the user get the 200 points or not? Without additional conditions, you cannot answer that. So you most likely need to introduce idempotency keys. That is a unique key that the client adds to each operation. You then check whether you have already processed that operation or not. If you have processed it, you do nothing. You just say okay, I already did it. If you haven't done it, you execute it.

Interviewer 2: There is also a question about partial release of points. When we hold some amount of points, but then some orders are canceled and some are completed, we need to partially release the hold. Is that possible?

Customer Expert: No, there are two options. The client wants to either fully release what they held or fully confirm it. There is no partial release. Either full confirmation or full release. Even if such a requirement existed, it would be better to implement it as a two-step process: first release everything, then hold again. But there is no such requirement, so it's all or nothing. However, as far as I remember, we had a rule that after 24 hours your system automatically releases the hold.

Interviewer 2: We also thought about making this period configurable in the request. For example, a bank might wait at most an hour for a response. If it doesn't come, something went wrong. But for an order, it could be a couple of days.

Customer Expert: Yes, that is probably a good API improvement. I think it would be useful.

Interviewer 2: Does anyone else from my team have questions? Apparently not. We also were asked to discuss user stories with you. What users might say. But we were told to write the user stories ourselves and discuss them among ourselves, then approve or reject them.

Customer Expert: Yes, look, user stories, despite having the word "user", do not necessarily have to be about an end user. For you, the user is not the end user, but another service.

Interviewer 2: We have two users. Who is the second? The end user, who will use the application that uses your API.

Customer Expert: Yes, good, valid case. You are right. So, for the end user, you most likely only have the option to see how many points they have. Because you probably won't hold points from the client side. In the client interface, there is probably only a checkbox to use bonus points, which goes not to your service but to some payment processing or order processing service. That order processing service sees the checkbox is checked and first does a hold, then a debit operation, then confirms the hold. So the user story to view bonus points does relate to your service. The user story about holding is not for the end user but for the client service. When the service is about to process an order that should debit points, it needs to find out how many points are available, hold the required amount, then confirm. So the second story is to release the hold if the payment fails. The last user story should be about the system itself, about releasing holds after a certain number of hours depending on how the points were accrued. You can put that in a user story where two services are used.

Interviewer 2: We have some drafts, but not in English. What's the best way?

Customer Expert: I have worked in English-speaking companies several times, I hope this won't be a problem.

Interviewer 2: Can I read them out?

Customer Expert: No, no, I'll read in the background. You have "shoot" and "must" like in RFC, right?

Interviewer 2: Yes, they gave us a formula and told us to write according to that formula. They said Moscow method of ranking points.

Customer Expert: Look, it's hard to say that giving points to the end user is a lower priority than holding, for example. But I think you could say that purely nominally, showing only available points is a must, and showing points that will burn in 30 days is a shoot. Or maybe a may. Because it would be strange if the user tries to spend points without knowing how many they have. Or if you accrue points and they don't know about them, it doesn't solve any user problem. But in the minimal version of your product, you can show only available points, not show held or burned points. You can split the requirements into two: assign must or shoot to one, and may to the other.

Interviewer 3: May I clarify, are you talking about the first one? Regarding the first, I meant the admin who is monitoring everything. You said you want to get the number of points in circulation to know how much you owe to users. That is what I meant.

Customer Expert: Oh, then maybe that is also a must? No? Well, okay. Shoot is shoot. It doesn't seem hard required, does it?

Interviewer 3: Thank you for your suggestions, we will definitely add them.

Interviewer 2: So, as for administration – do you more or less understand how you will check the second requirement? No, ideally it is better to delegate that to the user. Look, you say that bonuses should not be lost – that is a good thought, but in this wording, do you understand how you would verify it? It is most likely impossible to completely avoid that. I think you have the option to clarify the requirements. You can say: as an admin, I don't want double accruals, for example. Or I want all records, conditionally, that all accruals that came – say, 200 OK from an API call – are transactionally recorded in the database. We don't write the word "transactionally" in the requirements, but we formulate it exactly as: for all calls that returned a success status, we will not lose data. We only say that we implement this requirement not by our own means but through the transactional guarantees of the database. That is at least somewhat verifiable. But if you simply say that we will never lose data under any failures – then I come to you as an auditor and say: what if your disks physically break? Will you meet this requirement? You will say that most likely not, that wasn't in your threat model or risk model. That's why it's better to refine such requirements into testable ones. Because when you come to the defense, they will ask you: how will you prove it? Can you show it? So I would refine this to a testable requirement.

Interviewer 2: Different lifetimes, okay, fair enough.

Customer Expert: And now there seems to be nothing here about the following: if you try to do two holds in parallel and you implement it poorly, you could end up in a situation. Suppose a user has 200 bonus points. You make two holds of 200 each. If you do this not very correctly from a transactional point of view, you could get a situation where you have answered both clients that you successfully held the points, because you saw a database state where the user still had 200 points. So it would also be good to specify that you don't want this to happen. We've discussed that. In general, the rest of the requirements seem valid. Now, do we have something here about write-offs? That points disappear when their expiration date passes. We see that the user sees when points will expire. But we don't say here that the system should actually write them off. So we also want to do that. Look, we want to state here that we want to release points after some time. That is, in your scenario – use cases as we said – they may not necessarily be end-user cases or admin cases. The service consumer can also be a user. You can say that if the user hasn't confirmed the hold, they can expect – or should expect – that you will release those points no later than we agreed: one hour after the hold period expires.

Interviewer 2: Also, points that have expired will be written off no later than one hour after their expiration.

Customer Expert: And was there anything else in the original requirements? There was something as nice-to-have. It all depends on how much you intend to implement it. But we talked about wanting to write off the points that expire earlier first. That is, if I have been awarded points that will live for a year and points that will live for one month – even if the one-month points were awarded later – I still want the ones with the shorter lifetime to be written off first. This is also good to put into a user story. You can formulate it more generally: "I want points to be written off in order." Actually, how you formulate it matters. You can say: "I want to write off in the order they expire" – then you play in favour of the user. After all, this is a loyalty system to some extent, doing that will be nice for users. Or you can say: "I play in favour of the company – I will write off the longest-living points first." You probably don't want that. But you can put this into a user story so that it is clear how your system should work, as you plan it to work.

Interviewer 2: Thanks. Probably no more questions. Anyone else about implementation or any other things?

Customer Expert: It seems everything is fine.

Interviewer 2: Yes, but we started doing it and ran into the fact that we didn't fully understand exactly what we need to calculate and store – more precisely, how exactly to do it. We started writing the structure and thinking about what we will store.

Customer Expert: Don't hesitate to write. I would say that I may not answer the same day or the same hour, but most likely I will answer within a day. So if you get stuck on the second day or a week later and you don't really understand where to go next, it's better to ask me so that you still have almost a week to move forward. So write, I will try to answer of course.

Interviewer 2: Then if there are no questions – thank you, see you next week.

Customer Expert: Thank you, have a good evening.

Interviewer 2: Thank you very much.
