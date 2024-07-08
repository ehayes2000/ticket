<script setup lang="ts"> 
  const props = defineProps<{
    date: Date,
    padding?: string,
  }>();
  const year: string = props.date.toLocaleDateString("default", {year: "numeric"})
  const month: string = props.date.toLocaleDateString("default", {month: "short"})
  const day: string = props.date.toLocaleDateString("default", {day: "2-digit"})
</script>

<template>
  <div class="date-wrapper" :style="{padding: padding ?? '1rem'}">
    <div class="bg">
      <div class="pretty-date">
        <div class="month"> {{ month }} </div>
        <div class="day"> {{ day }} </div>
        <div class="year"> {{ year }} </div>
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

  /* .date-wrapper { 
    padding: padd
    padding: 1rem;
  } */

  .day { 
    font-weight: bold;
  }
  .year { 
    font-style: italic;
    padding-bottom: .3rem;
    font-size: medium;
  }
  .pretty-date {
    color: white;
    border-radius: .1rem;
  }
</style>