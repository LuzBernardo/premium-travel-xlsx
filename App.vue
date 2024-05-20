<template>
  <div id="app">
    <h1>Text Formatter</h1>
    <textarea v-model="userData" placeholder="Enter your data"></textarea>
    <button @click="generateFile">Generate XLSX</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      userData: "",
    };
  },
  methods: {
    async generateFile() {
      const result = await this.$guark.call("generate_xlsx", { data: this.userData });
      if (result.ok) {
        const blob = new Blob([result.data], { type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.download = "output.xlsx";
        link.click();
      } else {
        alert("Failed to generate XLSX file.");
      }
    },
  },
};
</script>

<style>
/* Add your custom styles here */
</style>

