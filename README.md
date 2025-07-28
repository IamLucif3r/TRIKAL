# TRIKAL : Tactical Reconnaissance for Incident Knowledge, Analysis & Leaks

Tactical Reconnaissance for Incident Knowledge, Analysis & Leaks (TRIKAL) is an automated cybersecurity news hunter and notifier. It fetches the latest cybersecurity news from trusted RSS feeds, filters for relevance using keyword/tag rules, and then applies a local Large Language Model (LLM) to further analyze and score articles for industry impact or significance. Only the most important news is sent as a concise, visually rich notification to your Discord channel.

## Features

- Aggregates news from your chosen RSS sources.
- Filters and deduplicates headlines with customizable rules.
- Uses a private, local LLM via Ollama for advanced relevance scoring.
- Notifies your Discord channel with top cybersecurity news using Discord embeds.

## How to Use

1. Configure your RSS feeds in rss.yaml.

2. Set your Discord webhook URL.

3. Run TRIKAL (Docker or standalone Go binary).

4. Get high-signal cybersecurity news automatically posted to Discord.


TRIKAL: The smarter way to keep your team (or yourself) up-to-date with the most relevant cybersecurity news.

