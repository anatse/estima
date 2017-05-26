<template>
  <div>
    ProjectPage {{user}}
  </div>
</template>

<script>
  export default {
    name: 'ProjectPage',
    data() { return {}; },
    computed: {
      user() {
        return this.$store.getters.getUser;
      },
    },
  };
</script>
