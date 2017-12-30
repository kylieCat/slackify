# Slackify
### Update your Slack status with your currently playing song and artist from Spotify

## Requirements
 - A Slack token: https://get.slack.help/hc/en-us/articles/215770388-Create-and-regenerate-API-tokens
 - Create an app on Spotify to get a client secret and ID: https://developer.spotify.com/my-applications/

## Usage
  - Clone this repo
  - Create a new file in this directory and name it `.env`
  - Add to `.env` file:

    ```
    SLACK_TOKEN="YOUR SLACK TOKEN"
    SPOTIFY_SECRET="YOUR SECRET"
    SPOTIFY_ID="YOUR ID"
    ```
    
  - Run `docker-compose up`. A URL will be printed to the terminal. Visit that URL in a browser to authenticate your app to access your Spotify account.
  - If everything worked the container will request your currently playing track from SPotify and set it as your status once every 60 seconds.
  - Output should look like this:

    ```
    slackify_1  | Login URL: https://accounts.spotify.com/authorize?client_id=XXXX&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&response_type=code&scope=user-read-currently-playing+user-read-playback-state+user-modify-playback-state&state=abc123
    slackify_1  | No Trigger - The (Not So) Noble Purveyors Of The Third or Fourth Coming
    slackify_1  | Banner Pilot - Spanish Reds
    ```

  - You can also run it with `docker-compose up -d; docker-compose logs -f slackify` to run as a daemon. Once you auth with the URL from the logs you can exit and the container will continue to run.
  - `docker-compose kill` to stop it.

## Known Issues
 - It will crash if an ad plays.
