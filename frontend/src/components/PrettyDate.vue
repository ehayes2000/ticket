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
    <div class="bg">
      <div class="pretty-date">
        <div class="month"> {{ month }} </div>
        <div class="day"> {{ day }} </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
  .bg { 
    width: 100%;
    height: 100%;
    background-color: blue;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    border-radius: .2rem;
  }

  .date-wrapper { 
    padding: .75rem;
  }
  .day { 
    font-weight: bold;
  }
  .pretty-date {
    color: white;
    border-radius: .1rem;
  }
</style>