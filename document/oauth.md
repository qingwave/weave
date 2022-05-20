# OAuth

Weave support auth login.

## Github OAuth
OAuth with Github, more details see https://docs.github.com/cn/developers/apps/building-oauth-apps/authorizing-oauth-apps
1. Open Github developer settings https://github.com/settings/developers, click `OAuth Apps`
2. Create OAuth app, click `New OAuth App` button
3. Register app, contents as follows
```
  name: Weave
  Homepage URL: http://127.0.0.1:8081
  Authorization callback URL: http://127.0.0.1:8081/oauth
  ```
4. Get your `clientId` and `clientSecret`, set in `config/app.yaml` and `web/src/config.js`

