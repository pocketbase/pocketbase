# Security

**Keep in mind that PocketBase is a non-commercial open source project, maintained entirely on volunteer basis (there is no company or dedicated team behind it), and there are no bounties!**

If you discover a security vulnerability within PocketBase, please send an e-mail to **support at pocketbase.io** or submit a private [GitHub Security advisory](https://github.com/pocketbase/pocketbase/security/advisories).

I try to be as responsive as possible and usually address security reports within a day or two, but if you didn't receive a reply from me for more than 5 days it is very likely that your email was flagged and in that case please open a GitHub issue or discussion just mentioning that you found a vulnerability and want to report it so that I can see the notification and will try to contact you for more details.

In case the vulnerability is confirmed, within another couple days I'll try to submit a fix, GitHub security advisory and CVE with remediation steps and **minimal details** regarding the found exploit to minimize giving too much hints to malicious actors (you'll be credited both in the fix release notes and in the public report).

### Please:

- DO NOT use LLM as part of your report or email communication - it is extremely frustrating to spend an hour or more reading a wall of generated text, writing an elaborate reply and in the end to just receive another generic LLM prompt response in return.

- DO NOT reserve and publish MITRE CVE number on your own and try to communicate first privately the details to better understand how the code is being used and whether the supposed vulnerability can be actually exploited in any real practical scenarios. Otherwise you are risking needlessly causing scaremongering and annoyance for users that rely on security scanners as part of their CI/CD pipeline.

- Wait before publicly disclosing and sharing your vulnerability, **ideally at least 5 days after the fix**, to give enough time for users to patch their instances _(you are free to provide a PoC and as much details as you want in your own blog post/gist/etc.)_.

### Below is a common list of vulnerabilities that were previously reported but are not considered a security issue:

<details>
<summary><strong>Stored XSS</strong></summary>

This was discussed several times, both privately and [publicly](https://github.com/pocketbase/pocketbase/discussions/6694), but I remain on the opinion that this should be handled primarily on the client-side.

Modern browsers recently introduced a basic [`Sanitizer` interface](https://developer.mozilla.org/en-US/docs/Web/API/Sanitizer) that could help filtering HTML strings without external libraries.

Having also a default [Content Security Policy (CSP)](https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CSP) either as meta tag or response header is a good idea to minimize the risk of XSS.
</details>

<details>
<summary><strong>SQL injection in low level DB methods like <code>app.DeleteTable(dangerousName)</code></strong></summary>

This is working as intended and it is not an issue but it is a common report most likely found by LLM or some other automated tools that may have stumbled on the [NB! code comments](https://pkg.go.dev/github.com/pocketbase/pocketbase@master/core#BaseApp.DeleteTable), warning exactly against it.

Raw SQL statements, table and column names are not parameterized and they are vulnerable to SQL injection if used with untrusted input. This is working correctly and it is not an issue. The documentation as seen above already warns against it.

In recent PocketBase releases, many of the arguments of these methods were also prefixed with `dangerous*` to make it even more clear that they should be used with caution.
</details>

<details>
<summary><strong>Race conditions</strong></summary>

To avoid DB locks PocketBase deliberately tries to minimize the use of DB transactions.
This means that operation like record update doesn't wrap out of the box for example the `SELECT` and `UPDATE` SQL statements in a single transaction, and this can technically lead to a race condition if multiple users edit the same record.

This is an accepted tradeoff and for the majority of cases it has no security implications.

This also apply for the read and delete of MFA and OTP records but for those cases, since they operate in a security sensitive context, they have an extra short-lived duration that can be configured from the collection settings _(they have also system cron jobs that takes care for deleting forgotten/expired entries to prevent accumulating invalid records)_.

For the cases where transactions are really needed, users can utilize the [Batch Web API](https://pocketbase.io/docs/api-records/#batch-createupdateupsertdelete-records) or [create a transaction programmatically](https://pocketbase.io/docs/go-records/#transaction) _(you can also wrap an entire hook chain in a single transaction)_.
</details>

<details>
<summary><strong>List/Search side-channel attacks</strong></summary>

Over the years we've implemented several extra checks to minimize the risk of List/Search side-channel attacks (see especially [v0.32.0](https://github.com/pocketbase/pocketbase/blob/master/CHANGELOG.md#v0320)) but users need to be aware that technically all client-side filtered fields are subject to timing attacks _(whether they are practical or not is a different topic)_.

This is by design and it is accepted tradeoff between performance, security and usability.

If you are concerned about timing attacks and have security sensitive collection fields such as `secret`, `code`, `token`, etc. then the general recommendation is to mark them as "Hidden" in order to disallow client-side filtering for these fields (superusers can still access them and they can take part of your API rules).
</details>

<details>
<summary><strong>Connecting to a vulnerable OAuth2 provider</strong></summary>

Because PocketBase v0.23+ supports automatically uploading the OAuth2 avatar on user create _(need to be specified from the auth collection OAuth2 fields mapping)_ some security researchers raised a concern regarding a Blind SSRF but this implies that an attacker controls the OAuth2 vendor and this is a very serious assumption in the first place.

The entire OAuth2 flow relies that the application server (PocketBase) trusts the configured OAuth2 vendor.
If you suspect that an OAuth2 vendor is malicious and cannot be trusted then you MUST NOT use that OAuth2 vendor at all and you should report it.

If someone is able to tamper with the OAuth2 responses then the entire OAuth2 flow can be thrown out of the window because they will be practically able to authenticate as any of your existing users and the eventual avatar url probing request is the least of your problem.

_Nonetheless, in future PocketBase releases there will be [extra `localhost` domain like checks](https://github.com/orgs/pocketbase/projects/2/views/1?pane=issue&itemId=159545722) when assigning the OAuth2 avatar url to a file field that will further minimize the risk of internal network probing requests in case of a vulnerable OAuth2 provider._
</details>

<details>
<summary><strong><code>disintegration/imaging</code> CVE-2023-36308</strong></summary>

Just for the past month, due to some corporate security scanner 5 different people raised concerns over [CVE-2023-36308](https://nvd.nist.gov/vuln/detail/CVE-2023-36308) but this is not really a vulnerability, especially not in PocketBase.

[`disintegration/imaging`](https://github.com/disintegration/imaging) is a direct PocketBase dependency responsible for the thumbs generation.

First, a panic (similar to exception in other languages) is NOT a security issue and Go programs usually have to be written with that in mind. In PocketBase specifically all routes have auto panic-recover handling, no matter what the source of the panic is, so the worst case scenario would be an HTTP error response when attempting to access the thumb.

Second, the related issue that the CVE describes is probably caused by a bug in an outdated `golang.org/x/image` dependency listed in the `go.mod` of that package but PocketBase uses a newer patched version of it that is expected to take precedence.

Third, even if that issue is still available, with PocketBase it would have been triggerable ONLY if we supported TIFF thumbs generation but we don't. The supported thumbs formats at the moment are JPG, PNG, GIF (its first frame) and partially WebP (stored as PNG). All other images are served as it is, without any transformation.

In the future I may consider eventually replacing the library because it is no longer actively maintained but as of now it is working correctly and as expected for our use case.
</details>
