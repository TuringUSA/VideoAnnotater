<template>
  <div class="is-flex is-vcentered content has-text-centered is-justify-content-center is-flex-direction-column container	is-align-items-stretch "
       style="min-height: 100vh">
  <upload-modal :open="isModalOpened" @close="isModalOpened = false"/>
    <div class="box p-5 ">
      <h1>Find content from the video</h1>
      <p>The system is using machine learning to find parts of the videos that contains the searched content. <br>
        You can add new video by clicking <a href="#" @click="isModalOpened = true">here</a>.
      </p>

      <div>
        <form @submit.prevent="search">
          <b-field>

            <b-input placeholder="Enter a keyword that you would like to find"
                     size="is-medium"
                     expanded
                     :loading="isRequesting"
                     v-model="keyword"
                     required
                     icon="file-search">
            </b-input>
            <p class="control">
              <b-button type="is-primary is-medium" :loading="isRequesting" native-type="submit">Search</b-button>
            </p>
          </b-field>
        </form>
        <p class="mt-2" v-show="hasSearched">Found <strong>{{results.length}}</strong> results for your search</p>

      </div>

    </div>

    <div class="container columns is-flex is-flex-direction-column is-align-content-stretch is-align-items-stretch is-justify-content-start p-5 mx-0">
      <video-result v-for="video in results" :file-name="video.Filename" :key="video.UUID" :uuid="video.UUID" :link="video.SignedURL" :features="video.Features"/>
    </div>

  </div>
</template>

<script>

import VideoResult from "../components/VideoResult";
import UploadModal from "../components/UploadModal";
export default {
  name: 'IndexPage',
  data() {
    return {
      isRequesting: false,
      keyword: "",
      results: [],
      hasSearched: false,
      isModalOpened: false
    }
  },
  methods: {
    async search() {
      this.isRequesting = true;
      try {
        const response = await this.$axios.get("/videos", {
          params: {
            keyword: this.keyword,
            page: 1,
            perPage: 50
          }
        })
        this.results = response.data.Videos ? response.data.Videos : [];
        this.hasSearched = true;

      } catch (e) {

      } finally {
        this.isRequesting = false
      }
    }
  },
  components: {UploadModal, VideoResult},
}
</script>

<style lang="scss">
body {
  padding: 20px;
  background-color: #dcdcdc;
}
</style>
