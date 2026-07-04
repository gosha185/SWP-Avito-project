# Customer Review Transcript – Week 5

**Date:** 4 July 2026  
**Participants:** Industry Expert (Customer), Team Members
**Meeting purpose:** Sprint Review – demonstration of Sprint 3 results, User Acceptance Testing (UAT), and feedback on the team's development process.

---

The speakers are referred to as **Team Member** and **Expert** (customer).

---

**[00:29]**

**Team Member:**

Good afternoon.

**[00:34]**

**Expert:**

Hello. Let's take a look at the other students. If we take a look. So, how many of us should there be? 


**[00:36]**

**Team Member:**

Well, in general, five of the members of the team.

**[00:40]**

**Expert:**

Oh, well, then it's kind of... Everyone is. Alright then, let's get started.

**[00:51]**

**Team Member:**

Hello, excuse me, please. I can't hear what you are talking from computer.
So, let's start. It turns out that we have a sprint preview again today. And the goal of our sprint on this week was to complete the MVP version 2. Here. What changes have happened to us? It turns out that we corrected bugs.

Including the bug that we said last week, where we made the wrong mistake. And also, some more technical... In the tent, let's say, we corrected it. And added the workers. They will then tell you in detail, it turns out, back-enders. So. There is a point to show the sprint backlog, probably, right?

**[01:04]**

**Expert:**

Let's see what we planned, what happened.

**Team Member:**

Now it turns out that our sprint backlog was completed at the moment at 6 in the process of 5. And there are still a lot of these days, because some of them will be done after our meeting.

That's why they are still here. In principle, I can also add to the sprint preview that we made daily scrubs on this week. We are slowly moving in this direction, much better than last weeks.

And... So. Probably that's all. I can tell you what our tasks were this week. We made the UML diagrams, only three types of them, statistical, dynamic and deployment. We also made architectural decisions records on which Katya will tell.

I think we'll probably start with the V2 version. 

**[01:17]**

**Expert:**

Yeah, let's see.

**[01:22]**

**Team Member:**

I think I did not pick a clear person for that. I forgot to mark this role. Let's probably start with the UML diagram, Gosha.

**Team Member:**

Hello. Can you hear me?

**Expert:**

Yeah, I can hear you.

**[06:57]**

**Team Member:**

Can you see the mouse? Yeah, yeah, yeah. So, we asked to make three diagrams.

The first component, the second one, the sequence, the sequence is our third deployment. This component is quite simple in essence. We have a database, a service, some user, sorry, we do something between them.

We have three main layers. It's a stack layer, it comes with a request, we call it a service layer, then a storage layer, and then we transfer the database to the database. Either we make a mistake, or we send it to the database, or we just call it a storage layer.

And we added a worker, which will be a part of the database. It will be directly called either a service or a storage layer of the method and the database will be cleaned by the transparent batches.

**Expert:**

Okay.

**[09:47]**

**Team Member:**

In this diagram, in the sequence, I wrote it in general, well, this project will be one, but I wrote it in a step, just like the team can accept so far, in general. Well, the request comes to the SAP layer, then to the service layer in the next method, then there is a very complex calculation between a storage layer and a service layer, and between the storage layer and the database itself. Well, here they are basically there is a request for the database from it, because there is a request for the database and nothing returns, because it's hard to categorize.

**[09:49]**

**Expert:**

Well, look, I just advise you here, I can go further, of course, separate teams, separate sequence with diagrams are described, because it's just hard to get into context. That is, the usual sequence with diagrams should not be huge. That's why it is possible to make it more logical to get a balance separately or to halve with confirmation Well, on halving with confirmation from the measure, you can do it through such a basic sequence in the chart which is called an alternative.


There is an alt block, it can be shown that, actually, if you have confirmation, then one file is generated, if you have halving, then the third is generated, if you have a timeout, then the fourth is generated. Here, you will have more diagrams, but they will be smaller by themselves, each of them. Yes, to redo some specific scenario.

**Team Member:**

Well, and then back to without any problems, and the worker also cleans, calls methods of service and storage layer. And the third diorama is about the main deployment. Here, in general, everything is also simple, there is some end user, he asks the API user to the server, and he already sends it to us to some server, or to the server leader, or to the usual server to the request, he processes the server, we have it with the worker, that is, the worker's task is to clean only one server, which you can see.

In general, it's quite simple, too, the scheme.
Well, yes, as if everything is good here.

**[09:57]**

**Team Member:**

Let's go then and discuss architectural solutions. Katya, can you please explain?

**[10:32]**

**Team Member:**

Yes, I will now open them. In general, we divided architectural solutions into solutions from our own service, from my layer, from the top, as far as I understand, I have two of them. This is just the design of the workers, which remove the finished batches after a certain time from the table, they are taken to a separate server and run only on it, not on everyone at once, if briefly.
I can show the text as a formula.

**[10:34]**

**Expert:**

Let's see.

**Team Member:**

We have four workers, two on the table, one of them is removed from the table or from the The other is removed from the table of the data after a certain time. That's how it is.

**[10:56]**

**Expert:**

In general, it's okay. yes, well, yes, and you gave it to you, yes, yes, and it turns out the second thing that I thought.

**Expert:**

Did they give you an ADR template?

**[12:20]**

**Team Member:**

Yes.

**[13:16]**

**Expert:**

Okay then. 

**[13:26]**

**Team Member:**

The second I thought can be can be considered such an architectural solution, this is just the fact that we have the removal of the logical by status is separated from the removal of the table at a certain time. 

**[13:40]**

**Expert:**
yes, well, look, in fact, you have the most basic idea, I could just be on how you generally store data so to guarantee independence of a separate transaction because in the first realization that you had, there was no such thing and you changed and you came up with a new problem and this can also be described in the idea and the worst practice is good thank you but in general, this is what we looked at the ADR 001, but this is also a good idea, how to implement the background processes. 

**Team Member:**

so if other backend members can tell you about other ADR, you get a word.

**Expert:**

Apparently no one is ready to tell about other ADR. 
We can then see if there is something else or see what you want to take in the next moment.

**Team Member:**

In general, we still have to demonstrate our own service, how it works, the base is fixed, the reports are not in the app and so on and separately demonstrate the workers because we almost didn't have time to fill them in the main branch, they are now locally for me, I will need some time to show them the demonstration,Leilia looked like she wants to discuss something.

**Team Member:**

Can you explain demonstration of what? 

**Team Member:**
Demonstraton of workers.

**Team Member:**
I can now show the file that we made with team - the development process - and then go to the workers.

**Expert:**

Let's do like this.

**Team Member:**

So we were asked to make a file of the development process, we of course worked all the teams on it because the whole team works and each of their own area works here for action, it turns out in general, how everything goes, what we do, the fact that we have a protected branch, we need to check and so on, a couple of examples, it is possible to give a feedback on the graph, it can be improved and an explanation, so to speak, the graph, then how we go through the process with problems on github definition of dan, how the code of another person is checked, this is the base of immigration, we also have automatic tests on github on the request and on the code and simple tests already on the code itself.

**Expert:**


In the template of the task definition of dan, yes, that is, you have a new task when it is created, I am automatically connected to the definition of done.


**Team Member:**

I am not sure about this, but how should it be connected?

**Expert:**

You have a template in github, you can write there a checklist of the definition of dan, and you will have a new task when it is created, the definition of dan that you wrote yourself, in fact, to be presented as a checklist, well, and the definition of ready, too, if you have it, it will also be presented as a checklist, you just go through it and press the buttons, well, that is, in every task, his copy is sent, then you press the buttons within the task. 

**Team Member:**

I understand well, I think we can try to do it, thank you, now it turns out to be a worker.

**Team Member:**

Yes, a second. As you can see, I understand correctly, there should be a visual studio code, I will show it on the already executed commands in the terminal so that it is faster, they are just demonstrating, in fact, for the first time, we clean the database, create some test users, create a hold with a fast flow of the hold, which we will cancel, we will make it old so that the worker, who is cleaning the table of data, has already been able to remove it, we now have a timeout of 30 days, if I am not afraid to hit this table, here we cancel, we make a fixed batch, here you can see that there are actually holds, batches, well, or rather, one hold, the workers are launched, you can see that they find these batches and the holds and do not remove them, then, for some reason, the batch table does not immediately update, but below after the restart, I restart because the workers have a certain amount, and if restart the server, they automatically restart the database, here it is already visible that everything is cleaned.

**Expert:**

Can I ask strange questions about why you are doing it so strange, check through the introduction of the text and like them through the level?

**Team Member:**

To be honest, it was written with a deeepseek, it shows how it works, and I did not understand how exactly it does it, 

**Expert:**

But the theoretical level should work if you have units, the price is not like that, judging by the fact that you do the translation of the text, it should not be saved as a text, well, in short, try to just like it, it should work much faster.

**Team Member:**

Well, thank you, that's all I have, it seems now you need to show the remaining service.

**Team Member:**

Are you ready Stepan?

**Team Member:**

Yes. Everything is ready, just now became, to be honest, now I want to demonstrate so can you see so from the main one, yes, we corrected the error that came back, the server error came back, I corrected it, but now we can see and the balance was also divided so that when we pulled the balance, only the balance came back without the need to introduce some days there, which in fact do not belong to the balance, and separately there was added an added point for returning the current points to the specific days, in short, for each one of us, we will go here, well, some of which are more interesting in general,

**Expert:**

if nothing has changed much, then I will probably be able to get out of it now, 

**Team Member:**

Now then I will show it locally, the local swagger because for some reason on the githab we from yesterday evening we approve our PRs and they are not allowed to merge for some reason and it is difficult because of this, now I would throw it on the volume, 

**Expert:**

in general, it is on the side of the conflict

**Team Member:**

there are no conflicts and it is just written that the review is required, although two or three people have been approved

**Expert:** 

then check who has a default reviewer, most likely some kind of limitation

**Team Member:**

so we counted, we made a cold, so it turns out that the error came back when the confirmation of the final code, yes, it seems to me that with a repeat code and with a repeat confirmation, so we confirmed it and now we will try again now so what is not so so something goes not according to the plan, it will be necessary to look again, okay, I will look again then why it returns 500 again, here about the new. Now the balance returns only the balance here is available and the expression is 30 days, separately, the sum is 30 days, these are the main changes in the handles, but it is not clear why, of course, I look again, of course, why 500 is still coming back, well, yes, this is the point, well, let's check the repeated call, well, the repeated call can not be yes, the repeated conforms, the repeated cancels should not be, of course, 500, it is expected to be corrected, it turns out that they showed everything they wanted, I can, as a result of the process of changing the errors, we added our product version 2, in fact, ready thanks to the worker, three notes were made from the architecture decisions records were made in half, that is, in the process, a little more need to bring up to three files development just with a sign, it turns out after our meeting today, we still work a little on the scenario, how can you work with our product and, in principle, that's all, well, yes, that is, it is useful.
if you made the dors and the dors, throw them in the template so that you have automatic bugs here, it would be necessary to correct and well, as in the UML, I would be the message sequence chart which broke the second one, after all, on the message sequence chart, which according to a separate scenario and the last one is here is the dr on how you store records about the change of the balance, this thing is like

**Team Member:**

Well, thank you for your time and feedback it turns out the next week

**Team Member:**

you can also find out what it may be, we have a goal for the next sprint, 

**Team Member:**

Yes, thank you

**Expert:** 

yes, well, that is, look, are you planning to finish on the second MVP?

**Team Member:**

It seems like we can do more, because on the 20th there is still a course,


**Expert:** 

yes, I would at least still achieve the finalization, well, finalize the things that we discussed today, yes, that is, the preparation of the documentation, that is, to complete the documentation according to the project, it can be a goal in general. You can do something like that.

**Team Member:**

There was an idea of making history user-friendly adding search by operation.

**Expert:** 

You lnow it is more interesting to make time filters than just filters. Like to filter transactions between two dates because usually there are complaints about having operation and not gaining points. In such cases ypu go and ask for transaction history to see with ypur own eyes.

**Team Member:**

Also there was an idea with optimization.

**Expert:** 

That's a good idea/ To make it you need somehow to measure what you are optimazing. You can not say it became better. If you know how to describe current and after improvisation optimisation, then it is  a good idea. Otherwise it is not.

**Team Member:**

Yes, we were taught assyptotic speed, so we could have database with information and try something.

**Expert:** 

Yes, ypu can try it. In such set up it is a normal situation.

**Team Member:**

Then again, is anyone wants to ask question? Then that's all. Thank you.

**Expert:** 

Yes, good day.

**Team Member:**

Good bye
