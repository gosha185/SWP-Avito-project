# Customer Review Transcript – Week 4

**Date:** 27 June 2026  
**Participants:** Industry Expert (Customer), Team Members
**Meeting purpose:** Sprint Review – demonstration of Sprint 2 results, User Acceptance Testing (UAT), and feedback on
the team's development process.

---

The speakers are referred to as **Team Member** and **Expert** (customer).

---

**[00:29]**

**Team Member:**

Hello everyone. Can you hear me?

**[00:34]**

**Expert:**

Hello. Yes, we can hear you.

**[00:36]**

**Team Member:**

Great. I assume we're still waiting for [inaudible]?

**[00:40]**

**Team Member:**

Most likely not. One of our teammates is currently in an area with poor internet access, so they're unlikely to be able
to join.

**[00:47]**

**Team Member:**

Alright then, let's get started.

**[00:51]**

**Team Member:**

Before we begin, could we please ask for your permission to record the meeting and create a transcript?

**[00:56]**

**Expert:**

Of course. Can I just give permission once for all future meetings, or do you have to ask every time? Either option is
fine with me.

**[01:04]**

**Team Member:**

I'm actually not sure. Every assignment specifically tells us to ask for permission before each meeting.

**[01:11]**

**Expert:**

I see. In that case, yes—you have my permission to record the meeting and create a transcript.

**[01:16]**

**Team Member:**

Thank you very much.

**[01:17]**

**Expert:**

You can treat this as a public meeting. Feel free to use or quote anything I say if necessary.

**[01:22]**

**Team Member:**

Great. Today we're following a slightly different format. We've been asked to conduct a Sprint Review with you as our
customer. Our sprint lasts one week, and this time we were instructed to present its results directly to you.

**[01:39]**

**Expert:**

Sounds good.

**[01:41]**

**Team Member:**

Could someone please share the Sprint Backlog? The page isn't loading on my side.

**[01:50]**

**Team Member:**

Let's wait a moment while someone opens the Sprint Backlog. [inaudible]

**[02:05]**

**Team Member:**

Let's open the Sprint Backlog for Sprint 2.

**[02:21]**

**Team Member:**

To briefly summarize our work this week, we introduced **Daily Scrums** as you recommended. Every day we discuss what
we've completed and what we're planning to work on next.

We also reorganized our Sprint Board into four columns: **To Do**, **In Progress**, **Review**, and **Done**. This
structure makes it much easier to track the progress of each task.

We haven't started every task yet, so today's Sprint Review is happening slightly before the sprint officially ends.
However, most of our tasks have already been completed or are currently in progress, and overall the sprint is going
well.

I think that concludes the Sprint Review itself. We were also asked to conduct **User Acceptance Testing (UAT)** with
you.

**[03:11]**

**Expert:**

Before we move on, I'd like to explain how Sprint Reviews are usually conducted in large technology companies so you
have a better understanding of what stakeholders typically expect.

Since I come from a technical background, many technical details are obvious to me. However, when you work in a large
company, it's considered good practice to invite the actual business stakeholders to Sprint Reviews. Those meetings
often include sales managers, marketing specialists, business development managers, and other people who know the
business domain very well but aren't deeply familiar with the technical implementation.

Because of that, Sprint Reviews should be understandable not only to engineers but also to non-technical stakeholders.

Of course, the format depends on your audience. If your customer is another engineering team, using technical
terminology is perfectly acceptable. But when presenting to business stakeholders, your explanation needs to focus on
business value rather than implementation details.

There are a few common practices that make Sprint Reviews more useful.

First, teams usually explain what their Sprint Goal was and how close they came to achieving it. This doesn't have to be
a long presentation. You simply say something like: "Our Sprint Goal consisted of implementing these two User Stories.
We fully completed the first one, partially completed the second one, and achieved about 50% of the planned work. We'll
explain why in a moment."

The second thing business stakeholders care about is what has actually been completed. That's where you demonstrate the
working functionality—which is exactly what we're about to do.

They also want to know what the team plans to accomplish during the next sprint. So the traditional Sprint Review
usually consists of three parts:

* what the previous Sprint Goal was;

* how successfully the team achieved it;

* what the next Sprint Goal will be and how the team plans to accomplish it.

The important thing is to explain everything in a way that business stakeholders can understand.

For example, if you simply say, "We're doing refactoring," many business people won't understand why that matters.

But if you say, "Next sprint we'll complete the frontend for this feature," or "We'll fully deliver this part of the
functionality," then it's immediately clear to them what value they'll receive.

There's one more important Agile principle I'd like to mention. Agile is more a collection of recommendations than a
strict framework.

One of its core ideas is that every sprint should deliver a potentially usable increment—a piece of completed
functionality that someone can actually use.

That means your Sprint Goal should also reflect this idea.

For example, implementing the entire frontend without the backend doesn't produce a usable feature.

On the other hand, delivering a small feature end-to-end—with both frontend and backend working together—is a valuable
increment because users can already benefit from it and solve at least part of their problem.

Hopefully that gives you a better picture of how Sprint Reviews are typically conducted in real product teams.

**[06:57]**

**Team Member:**

Yes, that explanation was very clear.

We also had a Sprint Goal for Sprint 2, although I honestly don't remember the exact wording.

To summarize, during the previous sprint we built the product, but it wasn't fully operational. Some components were
missing, and several tests were failing.

This week's objective was to implement the missing functionality, ensure the product passed all tests, and make
everything work correctly.

As far as I understand, we achieved that goal. Our tester spent quite a bit of time working on it.

Would anyone like to add something about the testing?

**[07:44]**

**Team Member:**

About the tests specifically?

**[07:46]**

**Team Member:**

Yes, about the product in general.

**[07:48]**

**Team Member:**

I implemented tests for concurrent requests and added most of the unit tests. I believe we've covered the majority of
the functionality.

**[07:59]**

**Team Member:**

Great.

As part of this sprint, we were also asked to conduct **User Acceptance Testing (UAT)**, so we'd like to demonstrate the
system to you.

Could someone please share their screen?

**[08:10]**

**Team Member:**

I think that's my responsibility. Since we don't have a graphical user interface yet, we'll demonstrate the
functionality through Swagger.

**[08:18]**

**Expert:**

Exactly. Since there isn't a UI yet, Swagger is perfectly fine.

**[08:23]**

**Team Member:**

Let's begin by awarding **500 bonus points** to a user. These points will expire in **30 days**.

The request completed successfully.

Now let's check the user's balance.

Here we can see that the user has **500 points**, all of which expire in **30 days**.

If we request the balance for points expiring within **10 days**, the result shows that **zero points** expire during
that period, which is exactly what we'd expect.

**[09:07]**

**Team Member:**

Now let's place **200 points** on hold for **24 hours**.

We'll also specify an **idempotency key**.

The request completed successfully.

We can now check the user's active holds.

As expected, **200 points are currently on hold**.

**[09:43]**

**Expert:**

Right.

What does the balance show now?

**[09:47]**

**Team Member:**

Here it is.

**[09:49]**

**Expert:**

That looks correct.

We placed **200 points on hold**, so there should now be **300 available points** remaining.

**[09:57]**

**Team Member:**

Exactly.

Now let's confirm the hold.

The hold was created using an order identifier, and it can later be either **confirmed** or **cancelled**.

We'll confirm it now.

The operation completed successfully.

If we check the list of active holds again, it should now be empty.

Yes—there are **no active holds**.

**[10:32]**

**Expert:**

And the balance?

**[10:34]**

**Team Member:**

Let's check.

**[10:36]**

**Expert:**

It should remain at **300 points**.

**[10:40]**

**Team Member:**

That's right.

Since the points were deducted from the available balance when the hold was created, confirming the hold doesn't change
the available balance any further.

That's why the balance remains **300 points**.

**[10:56]**

**Expert:**

Correct. That's exactly how it should work.

**[10:59]**

**Team Member:**

Now let's place another hold.

**[11:07]**

**Team Member:**

The request failed because we accidentally reused the same order identifier.

**[11:12]**

**Expert:**

That's actually an interesting example.

In practice, business requirements often change over time.

Imagine that customers receive **10% cashback** for an order, with those points expiring after three months.

Later, the business introduces a promotion that awards **another 10%** for the same order, but with a different
expiration period.

In that case, the same order could legitimately be associated with multiple bonus point batches.

So it's not always obvious that an order identifier should remain unique forever.

Your implementation is perfectly reasonable for the current requirements.

If business requirements change in the future, you can always revisit this decision.

I'm not suggesting that you change anything now—I'm simply illustrating how requirements evolve.

One important lesson is that you should never assume today's business rules are permanent.

**[12:20]**

**Team Member:**

Understood.

We've now successfully created another **200-point hold**.

Let's verify that.

Yes, we now have **200 points on hold**.

Next, let's cancel that hold.

We'll submit the cancellation request using a new idempotency key.

The hold has been cancelled successfully.

If we check the holds again, there should now be **zero active holds**, and the available balance should return to **300
points**.

That's exactly what happened.

We held **200 points**, cancelled the hold, and those points became available again.

**[13:16]**

**Expert:**

What happens if we try to cancel the same hold again?

Not with the same idempotency key—a completely new request.

**[13:26]**

**Team Member:**

Let's try it.

The operation fails, but currently the API returns **HTTP 500 Internal Server Error**.

This definitely needs improvement.

**[13:40]**

**Expert:**

Yes.

In this situation I'd expect **HTTP 400 Bad Request** instead.

The request itself is invalid rather than the server encountering an internal error.

**[13:52]**

**Team Member:**

Agreed.

We can also display the user's transaction history.

Here we can see the complete sequence of operations:

- bonus points awarded;
- points placed on hold;
- hold confirmed;
- another hold created;
- hold cancelled.

**[14:24]**

**Expert:**

I don't remember whether we discussed this previously, but did we decide that points with the **earliest expiration date
** should always be spent first?

**[14:34]**

**Team Member:**

Yes, I believe that requirement was included, and we've implemented it.

**[14:40]**

**Expert:**

Great.

Then next time we should verify that behavior.

For example, we could award points that expire in **30 days**, then award another batch expiring in **7 days**.

After creating a hold, we should verify that the system reserves—and later deducts—the points with the **shorter
remaining lifetime** first.

That behavior should also be visible when checking balances with different expiration filters.

I'd also like to test another scenario.

Let's award the same user **200 points twice**, using different idempotency keys for each request.

**[15:32]**

**Team Member:**

Alright.

Here's the first request—**200 points awarded**.

And here's the second one—another **200 points**.

Previously the user had **300 available points**, so the balance should now become **700 points**.

**[15:58]**

**Expert:**

Now let's try placing **all available points on hold**.

In other words, let's verify that your implementation can reserve points originating from multiple award transactions.

**[16:08]**

**Team Member:**

From multiple point awards, yes.

Let's create the hold.

Ah, the request failed because we accidentally reused the same order identifier.

**[16:21]**

**Expert:**

Right.

Please verify that scenario later.

**[16:25]**

**Team Member:**

Actually, the user has **700 points**, not 800.

I forgot that only **300 points** remained after the previous operations.

So let's place **700 points** on hold instead.

**[16:39]**

**Expert:**

Yes, let's reduce the amount accordingly.

**[16:42]**

**Team Member:**

Done.

**[16:44]**

**Expert:**

Great.

**[16:46]**

**Team Member:**

If we open the list of active holds, we can see that all **700 points** have been successfully reserved.

**[16:53]**

**Expert:**

Excellent.

The correctness of this scenario depends on your implementation details.

If your system explicitly links deduction transactions to individual award transactions, you'll need to maintain those
relationships carefully.

If, instead, your implementation simply operates on the total available balance, then this scenario becomes much
simpler.

Since everything works correctly in your implementation, I don't see any issues here.

**[17:18]**

**Team Member:**

I think that's everything we wanted to demonstrate.

Unless anyone has anything else to add, we can finish this part.

**[17:27]**

**Team Member:**

We'd mainly appreciate your feedback on today's demonstration.

**[17:34]**

**Expert:**

Before we continue, I noticed the browser icon on your desktop.

I actually worked on that browser for about a year and a half.

**[17:46]**

**Team Member:**

Really?

**[17:48]**

**Expert:**

Yes.

Small world.

Anyway, let's get back to the review.

**[17:57]**

**Team Member:**

Team Member, you wanted to ask a question.

**[18:02]**

**Team Member:**

Yes.

I have a question regarding our API responses.

Right now, after mutation requests, we return the internally generated **batch identifier**.

Is that actually useful?

More generally, what information should an API return after a successful mutation request?

**[18:28]**

**Expert:**

That's a very good question.

In most cases, the client is primarily interested in whether the operation succeeded.

If an error occurs, the response should clearly explain the reason.

If the request succeeds, it's often difficult to imagine that the client has forgotten what it originally sent.

One possible approach is to return the same fields that were provided in the request as confirmation that they've been
successfully stored.

Another perfectly acceptable option is to return an **empty response body** for successful operations.

There's nothing inherently wrong with that.

Error responses, however, should always contain meaningful information.

For example, earlier your API correctly reported **insufficient funds**.

Responses like that are genuinely useful to clients.

Personally, I think returning an empty response for successful requests is completely reasonable.

**[19:24]**

**Team Member:**

Thank you.

**[19:26]**

**Team Member:**

If there are no more questions, I think that's everything from our side.

We'd mainly like to hear your feedback on today's Sprint Review.

As for the next sprint, as you mentioned earlier, we still need to define our Sprint Goal.

At the moment our plan is simply to continue improving the project.

We'll probably receive additional university assignments, and those will naturally influence our priorities.

**[19:55]**

**Expert:**

Yes, and I'd like to emphasize one important point.

A **Sprint Goal** should never simply repeat the Sprint Backlog.

Its purpose is to identify what is **most important** during the sprint.

Why is that valuable?

Imagine you're halfway through the sprint and realize that completing every planned task is no longer realistic.

During your Daily Scrum, you should discuss how to redistribute your efforts so that you achieve the **Sprint Goal**
first.

The Sprint Goal is essentially the central objective of the sprint.

Everything directly related to that objective should receive the highest priority.

Other work is secondary.

That's exactly why the Sprint Goal exists—it helps the team make informed decisions whenever priorities need to change.

Ideally, the Sprint Goal should be defined **before the sprint begins**.

Of course, your project is educational, so the situation is a little different.

At **[redacted]**, for example, I always know which projects are expected to increase revenue or improve customer
satisfaction.

If my team's objective for a sprint is to deliver a feature that directly contributes to one of those business goals,
then that becomes our Sprint Goal.

If we later realize we won't have enough time to finish everything, I'd rather postpone technical debt or a nice-to-have
improvement than fail to deliver the work that actually affects revenue or user satisfaction.

That's why it's so important to define the Sprint Goal in advance, especially in real industrial projects.

Our Technical Director has an interesting way of describing teams that don't have a Sprint Goal.

He calls them **samurai**.

They have no goal—only the journey.

It's much better to have a destination before you start.

**[21:42]**

**Team Member:**

That's a great point.

Thank you very much for your time today, and thank you again for rescheduling the meeting.

This discussion was extremely helpful for us.

We'll see you again next week.

Goodbye.

**[21:56]**

**Expert:**

Goodbye.

Have a great day, and good luck with the project.

---