# Customer Review Transcript – Week 3

**Date:** 18 June 2026  
**Participants:** Customer representative (Expert from Avito), Team Members (Stepan, Georgii, Leilia, Ekaterina, Ivan)  
**Meeting purpose:** Sprint Review – feedback on MVP v1 scope, Product Backlog, Sprint Backlog, and team processes.

---

The speakers are referred to as **Team Member** and **Expert** (customer).

> **Team Member:**  
> Good evening. Hello. Hello. Let's get started, I think. We don't have any questions about the project in principle. We can share with you what we were asked and what we've been working on this week. Yes, let's take a look. Great. So we've done two things: the Sprint Backlog and the Product Backlog. In the Sprint Backlog, everything is still in "To Do" status – we haven't had time to update that yet. There are tasks we plan to do. Some of them are shown here as User Stories that we created. And for each task we had to come up with acceptance criteria. So we have task #27 here. I think we probably won't read each acceptance criterion, will we? Or what do you think?

> **Expert:**  
> No, probably not here. Maybe I can tell you what happens in Big Tech – I've worked in different companies. I can share what is usually used when creating tasks, what teams pay attention to and what they don't. Of course, everything is very individual, and you have your own assignment.

> **Team Member:**  
> Yes, that works for us. Could I have a moment? I forgot to ask for your permission to record this meeting and have a transcript. Yes, let's record it. Should I start the recording or will you? As far as I know, we have someone designated to record. Thank you very much for the permission. I think it's better to set up recording just in case. Let's do it – it won't hurt. Yes, we can also use it later.

> **Expert:**  
> Alright. When we create tasks, it's important to understand that tasks can be different. Conceptually, looking at your tasks – tell me, were you able to estimate all your tasks in story points, or not really?

> **Team Member:**  
> We were able to estimate all of them, but if you ask how accurate, I suspect not very.

> **Expert:**  
> Is there a difference – are you more confident in some tasks and less confident in others?

> **Team Member:**  
> Yes.

> **Expert:**  
> Why is that?

> **Team Member:**  
> Because some of the tasks have already been done. So the colleague who did them knows how many story points they can assign. But the tasks we haven't done yet are harder to estimate until we understand the scope of work.

> **Expert:**  
> Honestly, even the done ones are hard to estimate. Look, the further you go in your project, the easier it will likely become to estimate tasks in the part that already exists. For example, extending existing code and functionality is almost always easier to estimate and implement than writing from scratch. Because there are far fewer unknown factors. Such tasks are purely implementation – they have little research component – so at some point they will start to be estimated well. Over time, you'll get better at estimating, either in this project or when you join a larger company – you'll reach a plateau where estimates match tasks in the domain you know well. But you will certainly have a second class of tasks – research tasks. These are tasks that involve investigating something or writing something where you don't yet have expertise. And there, your estimates will fluctuate almost always. First, they tend to have higher estimates – not 1–2 story points but rather 5, 8, 13. Second, the variance will be very large – you might estimate 13 and do it in 5, or estimate 5 and do it in 13. Honestly, there is no perfect way to deal with this, but there are two approaches that help manage it. The first is to take such a small research piece, a small scope, that it still provides value but doesn't consume an insane amount of time. For example, in many Big Tech teams, they try not to set estimates above 5 – they decompose tasks so that the maximum estimate per person per sprint is about 5 story points. Of course, it depends on what one story point means, but conceptually, avoid very large tasks. The second approach – which is obvious but many teams forget, and when they remember it becomes easier – is that you don't have to fit all research into the story points you assigned. You can do what's called a spike or time‑boxed research. You go from the opposite direction. You don't say, "I will do all the investigation needed to solve this task in 8 story points." Instead, you say, "We as a team decided to spend 8 story points on research. Whatever we get done in that time is what we deliver." So you evaluate not by the volume of work but by the time you're willing to spend. So there are several estimation approaches. For tasks you understand well and are in your domain, you'll eventually get decent predictions. For research tasks, two key things to remember: limit the time from above for tasks that are completely unclear, and don't take super‑large estimates.

> Now you already have tasks and have decomposed them. For tasks, first of all, there should be a description. I saw you have it – you say you need to implement a user story. But there's an interesting case. If the whole team is deeply familiar with the code, that might be enough. If not everyone is deeply familiar, it's often good for someone who knows the code to make a more detailed decomposition – for example, if you're working with code‑centric tools, you might have a plan‑mod. In general, if you're working with code, you may need less of that. But if you're working with real people who aren't yet immersed in the code, it's more necessary. Besides the user story, it's good to describe at a high level what needs to be done: add such‑and‑such repository method, add such‑and‑such use case, add such‑and‑such handler, then assemble it all, wire it in the startup, and add tests – integration tests, unit tests – those kinds of things.

> Then you showed me that you have acceptance criteria for tasks – that's great. Acceptance criteria help you agree on what "done" means for a task. But there are two other things often used when creating tasks or verifying that you actually did them: Definition of Ready and Definition of Done for a task. What is Definition of Ready? It's your internal checklist for when a task can be taken into work – for example, all prerequisites are done, all ambiguities are resolved. In Big Tech, what usually goes into DoR? Things like designs are ready, the discovery document is finalised, you have business metrics to collect, you've written and signed off an architectural document. Only then can you take the task into work. Then there's Definition of Done – that's about when we consider the task truly closed and moved to final status. It's good to understand how DoD differs from acceptance criteria. In acceptance criteria you usually write a product story – not always, but typically you write things like "As a user I can do this", "As a user I can do that", and "when I do this action, I get an error". That's all fine, but you haven't said, for example, that the task is considered done when you've written unit tests for all methods and achieved at least, say, 65% coverage, or when you have at least one integration test for the new functionality. Or, for example, in DoD you can state that you consider the task closed when it has reached production. For instance, my team deploys 1 to 5‑7 times a day – not because of hotfixes, but because we have many services and we deploy constantly. And we have it explicitly in our DoD that a task is closed only when it's in production, not when a developer just says "I'm done". To make this convenient, if you work with GitHub, you can put templates – Markdown templates – so that every new issue automatically gets a checklist if you have standard DoD and DoR. That will apply to newly created issues – it won't retroactively apply to old ones, but it can for new ones. That can be handy. But I'm not insisting – I'm just telling you how it works in Big Tech. I hope that was useful.

> **Team Member:**  
> Yes, it was definitely useful. Thank you very much. For example, I saw acceptance criteria, and the template, and simple formatting, and there were also questions about how it differs from Definition of Done – thanks for explaining that.
> So MVP1 – the backend developers have started working on it, and tests too, as far as I know. Version v0 is already finished, and v1 will be conducted when the version is ready. For now, we don't have anything to show you at this meeting – we just set up the board.

> **Expert:**  
> Look, this backlog – usually we put the whole list of tasks we have for the nearest period. Typically we form a sprint from the backlog and run the sprint with its own board. Do you have something like that?

> **Team Member:**  
> Absolutely correct, yes. Here we have all the issues, but these two are not included in our sprint because we gave them a "should have" priority. All the rest are in a separate project for the sprint, and there we can mark task completion.

> **Expert:**  
> I haven't worked with GitHub for about 4 years, but I think you can set up a board similar to Kanban – not Kanban exactly, but a board like in Trello, where you drag cards and each day you can see how tasks move across columns. That visualises progress quite well.

> **Team Member:**  
> Okay, I'll do that. Thank you. I haven't read much but I couldn't figure out how to set it up – I'll look into it and do it. So that will be our Sprint Backlog, right?

> **Expert:**  
> Yes, that will be the sprint backlog. And if everything is estimated – and it seems everything is – then over time you should have a burndown chart. That's the second thing that's very useful to look at. You have columns where you move tasks, and you can see how the sprint's remaining work changes over the sprint days – whether you've moved tasks to "ready to deploy" or "done". The second thing that helps you understand how you're tracking against plan. It's not about judging good or bad – it's about understanding whether what you planned matches reality. For example, if you take 100 story points into a sprint and complete 80, and then in the next sprint you take 100 and complete 80 again, you might start taking 80 story points. Two data points isn't a large sample for statistics, but conceptually your intuition should tell you that if twice in a row your team's capacity is honestly 80, aim for 80. The burndown chart shows days in the sprint with an ideal burn‑down line. You don't need to achieve a perfectly smooth burn, but you should see that your actual burn follows that line with some variance. As far as I remember, GitHub can build that out of the box, but they've added enterprise features – maybe it's now in the enterprise tier, and if so, never mind. If it's still in the free version, I would set it up because it helps the team visualise how well you meet your own plans.

> **Team Member:**  
> Okay, thanks. So, any other questions? I'm sure someone has questions. I think we also have a requirement to specifically learn about the implementation of MVP1 – what functionality will be included in MVP1.

> **Expert:**  
> Let's see what you've potentially planned for MVP1. Conceptually, it seems ideal to reach a state where the service no longer has the problems we identified – all problems are fixed. And ideally, in each sprint we want to tangibly move closer to that goal. What does that mean? If we solve at least one or two problems per sprint, and we have a finished increment – a service where the problem previously existed but no longer does, and we can demonstrate that somehow via tests or otherwise – we can consider that a good increment. So you need to understand how much you can get done, and your goal in these sprints (weekly or bi‑weekly) is to show some finished increment. If you close one problem, that's already a finished increment – good. Close two – even better. But what would be not so good is if you said, "We solved each of the six problems (or however many) by 20%." That's not great because in a large company with a big business, solving six problems by 20% isn't very interesting. They'd much prefer you honestly solved one, closed it, it's off their plate, and they can breathe a sigh of relief. Even if solving each problem sequentially takes longer than doing them all in parallel, businesses usually still agree to incremental improvements because they understand what you're doing and see progress. So my point is: if in the sprint you take tasks that show the service has become incrementally better, and it's a finished increment, then it's okay with me.

> **Team Member:**  
> So for MVP, it's better to just close what was stated in the requirements? As the initial part, yes.

> **Expert:**  
> Alright. Any more questions? If not, I'll throw in another thought – not a task, just a thought. You're simulating a real work process. Besides the backlog and working on it, you should have other Agile rituals. One of them is the retrospective, which you'll likely have after you close the first sprint. My advice: write down the problems you encounter somewhere – in a text file, for example – because if you don't, at the retrospective you'll only remember what went wrong on the last day. In reality, if you write them down, you'll have much more material to discuss and want to discuss. Second thing: if you do have a retrospective (I hope you will, because it's good for improving processes), remember there are two things to discuss. First, metrics – you'll understand how many story points you planned versus how many you actually did. That's real metrics. Second, not metrics but processes – your internal feelings about how well the project is going, or not. We usually discuss those within teams. So we separate trying to improve processes based on metrics, and trying to analyse our internal grumbles about things we don't like and improve that too. So my advice: don't delay – write down problems you encounter early, before the retro. And when you run the retro, remember there's a part about metrics, which is objectively assessable, and a part about overall feelings – what's going well and what's not. Write those down too.

> **Team Member:**  
> Yes, we do have that – the retrospective, or something like that. Thanks for explaining what to put in it.

> **Expert:**  
> And the person facilitating the retrospective would do well to prepare the metrics from the sprint themselves. So you come to the retro – everyone has their own questions and suggestions, but the facilitator usually also brings some factual metrics. Start simple – the volume you completed versus planned. Later you can calculate more complex metrics. But let's do your first retro, come back with questions, and I might suggest what else you can measure.

> **Team Member:**  
> Okay, thank you very much. Well, if there are no more questions, we can wrap up. Again, please don't hesitate to write if any questions come up – I'll try to respond promptly. Thank you very much. Have a good evening. Until next time.

> **Expert:**  
> Yes, and have a good evening too. Good luck. Thank you, goodbye.
