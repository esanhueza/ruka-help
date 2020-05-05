<template>
  <div id="board">
    <div id="tabs">
        <div v-bind:code="this.code" v-on:click="onTabClick('code')" class="tab-item" :class="{'tab-item-selected': activeTab === 'code' }">
            <span>CODE</span>
        </div>
        <div v-on:click="onTabClick('schema')" class="tab-item" :class="{'tab-item-selected': activeTab === 'schema' }">
            <span>SCHEMA</span>
        </div>
        <div v-on:click="onTabClick('data')" class="tab-item" :class="{'tab-item-selected': activeTab === 'data' }">
            <span>DATA</span>
        </div>
        <div v-on:click="onTabClick('attachments')" class="tab-item" :class="{'tab-item-selected': activeTab === 'attachments' }">
            <span>ATTACHMENTS</span>
        </div>
    </div>
    <div id="board-content-wrapper">
        <CodeBoard v-if="activeTab   === 'code'"></CodeBoard>
        <SchemaBoard v-if="activeTab === 'schema'"></SchemaBoard>
        <DataBoard v-if="activeTab   === 'data'"></DataBoard>
        <AttachmentsBoard v-if="activeTab === 'attachments'"></AttachmentsBoard>
    </div>
  </div>
</template>

<script>

import CodeBoard        from '../components/boards/CodeBoard.vue'
import SchemaBoard      from '../components/boards/SchemaBoard.vue'
import DataBoard        from '../components/boards/DataBoard.vue'
import AttachmentsBoard from '../components/boards/CodeBoard.vue'

export default {
    name: 'Board',
    components: {
        CodeBoard,
        SchemaBoard,
        DataBoard,
        AttachmentsBoard
    },
    props:[
        'code',
    ],
    data: function(){
        return {
            activeTab: "code",
        };
    },
    methods: {
        onTabClick: function (tabName) {
            this.activeTab = tabName;
        }
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
#board{
}
#board-content-wrapper{
    height: calc(100% - 20px - 30px);
    width: 100%;
    border-top: 1px solid black;
    background-color: #c3c3c3;
}
#tabs{
    margin-top: 20px;
    height: 30px;
    display: flex;
    flex-direction: row;
    justify-content: space-around;

}
.tab-item>span{
    display: inline-block;
    line-height: 30px; 
    vertical-align: middle; 
}
.tab-item{
    height: 30px;
    width: 160px;
    border: 0px solid black;
}
.tab-item.tab-item-selected{
    border: 1px solid black;
    border-bottom: 0px solid white;
    background-color: #c3c3c3;
}
</style>
