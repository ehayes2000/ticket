<script setup lang="ts"> 
  const props = defineProps<{
    date: Date
  }>();
  const year: string = props.date.toLocaleDateString("default", {year: "2-digit"})
  const month: string = props.date.toLocaleDateString("default", {month: "short"})
  const day: string = props.date.toLocaleDateString("default", {day: "2-digit"})
</script>

<template>
  <div class="date-wrapper">
    <div class="pretty-date">
      <div class="month"> {{ month }} </div>
      <div class="day"> {{ day }} </div>
    </div>
  </div>
</template>

<style scoped>
  .date-wrapper { 
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }
  .day { 
    font-weight: bold;
  }
  .pretty-date {
    background-color: blue;
    color: white;
    padding: 0 .5rem 0 .5rem;;
    border-radius: .1rem;
  }
</style>