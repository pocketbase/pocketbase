# pocketbase-sync
Video Demo: https://youtu.be/91_0x8StFl0
# Warning
Well, several warnings.
1. Aside from this readme document this entire project was nearly all vibe coded. Accept that however you wish.
2. I'm a pre-junior developer who is fascinated by system architecture and infrastructure such as this, but I built this in my spare time and may not be the best FOSS lead for a project. Do with that what you will.
3. All of my personal applications are too small to need horizontal scaling... I built this because it was fun, and I'm fascinated by the tech. So if you choose to use this, walk into it knowing that I'm basically a chef who doesn't eat his own food.
4. Backup your stuff! While I've done some testing and have had good results, I would certainly not call this production ready. If my program eats your data and you're without a backup, you are to blame.
5. NATS Jetstream, used here, is NOT encrypted. You will need to come up with your own security scheme so that your sync endpoints aren't exposed to the public internet. Tailscale and ZeroTier are both excellent choices, but there are others out there. (I think K8s & K3s have some sort of provision for this, which reminds me, Fly.io should have something you can use also)
6. Normal record syncing has been tested to some extent, but I haven't had a chance to build and test any kind of realtime application for the realtime websocket syncing.
7. I might have built this whole thing just to solve another one of the internet's great non-problems, you know, like bitcoin did for finance. Basically, this whole project may be useful to a total of 0.000 of people out there. Shucks.

Now that the warnings are out of the way....
# What's going on here?

pocketbase-sync is Pocketbase (amazing) but with horizontal sync capabilities.
 - I'm using NATS Jetstream for the pub/sub architecture
 - This is Event driven, so it won't eat up your system resources like polling solutions such as the very-cool project marmot which syncs sqlite databases (I'm grateful to marmot for inspiring this project!)
 - Leaderless - there is no single central node that all other nodes must talk to. Once a sync chain has nodes talking on it, you can destroy and start any node you wish.
 - Eventually Consistent (usually syncs in miliseconds, but we aren't strongly consistent)
 - Realtime database connections are handled by websocket routing, so that no matter what instance an end-user originates from... if a realtime connection is open to a resource, all other users will connect to that same resource. Whoever initiated the realtime connection is the winner, and will hold the lease until the instance dies or until there are no realtime subscribers. This has not been tested!!!
 - Snapshots of your database are configured for the NATS Jetstream so that when a new instance joins, it can initialize using a recent snapshot of your database, followed by all database changes since the snapshot. It wouldn't be reasonable for a new instance to have to ingest every change since the beginning of time when it joins. So it grabs and applies the snapshot first, then catches up to current and begins syncing itself.

# Usage
Well, the command line side of things is unchanged. It's just Pocketbase all the way down. I added `plugins/sync` and made some minor UI tweaks to support it.

In the Pocketbase Admin UI, navigate to Settings => Application => Sync.

For the first instance, check the button that says "Initial Node" and fill in the local address that you want NATS to listen on, as well as an instance ID so you can know which one this is.

For all other following instances that you want to sync, I strongly recommend starting with a brand new empty instance of Pocketbase-sync, and in the Sync settings use "Add To Sync Chain" and similarly, fill out the local instance ID and NATS URL, but also fill in the Remote credentials of another instance you have on the sync chain.

Since this is leaderless, it doesn't matter what instance you point a new one at. It should catch up with all changes.

<img src="https://f004.backblazeb2.com/file/mattzab-dropbox/2025-12-08_17-56.png">



# Contributions
 - It would be nice to have a realtime app I can use to test the realtime websocket routing. If anyone has something they want to volunteer, I'm open to that.
 - If anyone wants to look over the codebase (almost everything new is in plugins/sync) and make sure there's some sanity going on in there, that would be nice of you.


What else should I put in this readme file? I wrote this myself, sans-LLM. So that's why it reads a bit rough. But hey, I write like I talk. This is what you get.

