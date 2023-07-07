<template>
  <b-modal v-model="isOpened" scroll="keep">
    <div class="card p-5">
      <h1>Upload new video</h1>
      <p>New video will be uploaded and processed in single request, so please don't close the tab while viedo is processing.</p>

      <form @submit.prevent="upload">
      <b-field>
        <b-upload v-model="file" drag-drop expanded accept=".mp4" required :disabled="isRequesting">
          <section class="section">
            <div class="content has-text-centered">
              <p>
                <b-icon icon="upload" size="is-large"></b-icon>
              </p>
              <p>Drop your video file here or click to upload (only .mp4)</p>
            </div>
          </section>
        </b-upload>
      </b-field>
      <p><strong>{{file.name}}</strong></p>

      <b-button type="is-primary" expanded :loading="isRequesting" native-type="submit">Submit</b-button>
      </form>
    </div>
  </b-modal>
</template>

<script>
export default {
  name: "UploadModal",
  data(){
    return {
      isRequesting: false,
      file: {}
    }
  },
  props: {
    open: {
      type: Boolean
    }
  },
  computed: {
    isOpened: {
      get(){
        return this.open
      },
      set(value){
        this.$emit("close")
      }
    }
  },
  methods: {
    async upload(){
      try {
        this.isRequesting = true;
        let formData = new FormData();
        formData.append("video", this.file);
        await this.$axios.post("/videos", formData)
        this.$parent.$buefy.toast.open({
          duration: 5000,
          message: 'Video uploaded successfully!',
          position: 'is-top-right',

          type: 'is-success'
        })

        this.isOpened = false
      }catch (e) {
        this.$buefy.toast.open({
          duration: 5000,
          message: `Something went wrong`,
          position: 'is-top-right',
          type: 'is-danger'
        })
      }finally {
        this.isRequesting = false;
      }
    },
  }
}
</script>

<style scoped>

</style>
