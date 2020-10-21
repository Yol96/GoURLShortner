<template>
  <div class="app-generate-form">
    <div class="logo-header">
      <h1 class="display-5">
        <img src="@/assets/logo.png">
        GoUrlShortner
      </h1>
    </div>

    <form class="mt-4" @submit.prevent="generateUrl">
      <div class="form-group">  
        <div class="input-group">
          <div class="input-group-prepend">
            <span
              class="input-group-text"
              id="validationTooltipUsernamePrepend"
              >URL: </span
            >

          </div>
          <input
            class="form-control"
            v-model="formData.address"
            aria-describedby="emailHelp"
            placeholder="http://www.google.com/"
          />

        </div>
      </div>
      
      <div class="form-group">
        <div class="input-group">
          <div class="input-group-prepend">
            <span
              class="input-group-text"
              id="validationTooltipUsernamePrepend"
              >Minutes:</span
            >

          </div>
          <input
            type="number"
            class="form-control"
            v-model="formData.expiration_time"
            aria-describedby="emailHelp"
            placeholder="10"
            min="0"
            max="60"
          />

          </div>
      </div>

      <div class="form-group">
        <button class="btn btn-primary" href="#" role="button">
          Generate
        </button>
      </div>
    </form>

      <hr class="my-4" />
      <p class="lead" v-if="shortUrlRecieved">
        Generated URL:
        <a target="_" :href="`${baseURL}${responseData.short_link}`">{{
          `${baseURL}${responseData.short_link}`
        }}</a>
        <br />
        {{
          `Created at: ${responseData.created_at}`
        }}
      </p>
  </div>
</template>

<script>
export default {
  name: "MainForm",
  components: {},
  data() {
    return {
      shortUrlRecieved: false,
      formData: {
        address: null,
        expiration_time: null
      },
      responseData: {
        address: "http://google.com",
        expiration_time: 0,
        created_at: "2020-10-19 13:30:53.21417948 +0300 MSK m=+517.111964179",
        short_link: "a"
      },
      err: null
    };
  },
  methods: {
    generateUrl() {
      this.shortUrlRecieved = false;
      if (
        this.formData.address == null ||
        this.formData.address == ""
      ) {
        this.shortUrlRecieved = false;
      } else {
        let finalFormData = {
          address: this.formData.address,
          expiration_time: this.formData.expiration_time
        };
        this.axios
          .post("/new", finalFormData)
          .then(res => {
            this.responseData = res.data;
            this.shortUrlRecieved = true;
          })
          .catch(err => {
            this.shortUrlRecieved = false;
            this.err = err;
          });
      }
    }
  }
};
</script>

<style lang="scss" scoped>
.app-generate-form {
  width: 30%;
	position: absolute;
	top: 50%;
	left: 50%;
	transform: translate(-50%, -50%);
  color: #eeeeee;
  .form-group {
    width: 100%;
	position: relative;
    .input-group-text {
      width: 6em;
    }
  }
}
</style>