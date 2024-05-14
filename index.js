require("dotenv").config({ path: __dirname + "/.env" });
const { twitterClient } = require("./twitterClient.js")
const axios = require('axios');

const tweet = async () => {
  try {
    const tweetText = `The URL **${storedSiteUrl}** has a safety status of **${storedSiteSafety}** on **Netsepio**`;
    await twitterClient.v2.tweet(tweetText);
      } catch (e) {
    console.log(e)
  }
}

let storedSiteSafety = "";
let storedSiteUrl = "";

const getLatestReview = async () => {
    try {
        console.log("stored", storedSiteSafety, storedSiteUrl)
        const response = await axios.get(`https://gateway.netsepio.com/api/v1.0/getreviews?page=1&domain=`);
        const review = response.data.payload.reviews[0];
        const { siteSafety, siteUrl } = review;

        if (siteSafety === storedSiteSafety && siteUrl === storedSiteUrl) {
            console.log("No updates");

        } else {

            storedSiteSafety = siteSafety;
            storedSiteUrl = siteUrl;
            
            console.log("Updated values:");
            console.log("Site Safety:", siteSafety);
            console.log("Site URL:", siteUrl);

            tweet()
        }
    } catch (error) {
        console.error("Error fetching data:", error.message);
    }
};

const runEvery30Seconds = () => {
    setInterval(async () => {
        await getLatestReview();
    }, 30000);
};

runEvery30Seconds();

