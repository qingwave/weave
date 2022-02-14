<template>
  <div class="flex flex-wrap w-full h-full p-2rem pb-5rem bg-gray-100 justify-between">
    <div class="w-32/100 h-10rem bg-white rounded-md mb-1rem" v-for="(item) in cards" :key="item.name">
      <div class="flex w-full h-full px-3rem items-center">
        <div class="flex-col w-6/10">
          <div class="text-gray-500">{{ item.name }}</div>
          <div class="text-4xl text-bold my-0.5rem">{{ item.value }}</div>
          <div class="inline-flex items-center">
              <arrow-circle-up v-if="item.up" theme="filled" size="18" fill="#4ade80"/>
              <arrow-circle-down v-else theme="filled" size="18" fill="#f87171"/>
              <span class="text-gray-400">{{ item.delta }}</span>
          </div>
        </div>
        <div class="w-4/10  text-center">
          <component :is="item.icon" theme="outline" size="48" :fill="item.iconColor" />
        </div>
      </div>
    </div>

    <div class="w-49/100 h-20rem bg-white rounded-md mb-1rem">
      <v-chart class="w-full h-full" :option="gaugeOption" />
     </div>

     <div class="w-49/100 h-20rem bg-white rounded-md mb-1rem">
      <v-chart class="w-full h-full" :option="pieOption" />
     </div>

    <div class="w-full h-20rem bg-white rounded-md mb-1rem">
      <v-chart class="w-full h-full" :option="lineOption" />
     </div>
  </div>
</template>

<style scoped>
</style>

<script setup>
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { PieChart, GaugeChart, LineChart } from "echarts/charts";
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from "echarts/components";
import VChart from "vue-echarts";
import { ref } from "vue";
import { Peoples, ArrowCircleUp, ArrowCircleDown } from '@icon-park/vue-next'
import { IconPark } from '@icon-park/vue-next/es/all';

use([
  CanvasRenderer,
  PieChart,
  GaugeChart,
  LineChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
]);

const cards = ref([
  {
    name: 'PV',
    value: 31284,
    up: true,
    delta: '10.5%',
    icon: 'book-open',
    iconColor: '#bbf7d0'
  },
  {
    name: 'UV',
    value: 21478,
    up: false,
    delta: '10.5%',
    icon: 'peoples',
    iconColor: '#fed7aa'
  },
  {
    name: 'APP',
    value: 125,
    up: true,
    delta: '10.5%',
    icon: 'application-one',
    iconColor: '#bae6fd'
  },
]);

const lineOption = ref({
  legend: {
    bottom: '5%'
  },
  xAxis: {
    type: 'category',
    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      data: [150, 470, 392, 218, 135, 248, 290],
      type: 'line',
      name: 'This Week' 
    },
    {
      data: [234, 360, 410, 326, 291, 310, 341],
      type: 'line',
      name: 'Last Week' 
    },
  ]
});

const pieOption = ref({
  legend: {
    bottom: '5%',
    left: 'center'
  },
  series: [
    {
      name: "Traffic",
      type: "pie",
      radius: ['60%', '70%'],
      avoidLabelOverlap: false,
      label: {
        show: false,
        position: 'center'
      },
      data: [
        { value: 335, name: "Direct" },
        { value: 310, name: "Email" },
        { value: 625, name: "API" },
      ],
      labelLine: {
        show: false
      },
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: "rgba(0, 0, 0, 0.5)"
        }
      }
    }
  ]
});

const gaugeOption = ref({
  series: [
    {
      name: 'Success',
      type: 'gauge',
      detail: {
        formatter: '{value}'
      },
      data: [
        {
          value: 50,
          name: 'SCORE'
        }
      ]
    }
  ]
});

</script>
