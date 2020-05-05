<template>
  <div id="board-content">
    <codemirror ref="myCm"
            v-model="code"
            :options="cmOptions" 
            @ready="onCmReady"
            @focus="onCmFocus"
            @input="onCmCodeChange">
    </codemirror>
  </div>
</template>

<script>

import { codemirror } from 'vue-codemirror'

import 'codemirror/lib/codemirror.css'
import 'codemirror/mode/sql/sql.js'

export default {
    name: 'SchemaBoard',
    components: {
        codemirror
    },
    data: function(){
        return {
            code: 'CREATE TABLE',
            cmOptions: {
                // codemirror options
                tabSize: 1,
                mode: 'text/x-sql',
                theme: 'default',
                lineNumbers: true,
                style: {
                    height: "100%",
                },
                line: true,
                viewportMargin: Infinity,
            }
        };
    },
    props: {
    },
    methods: {
        onCmReady(cm) {
            console.log('the editor is readied!', cm)
        },
        onCmFocus(cm) {
           console.log('the editor is focus!', cm)
        },
        onCmCodeChange(newCode) {
            console.log('this is new code', newCode)
            this.code = newCode
        }
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

#board-content{
    padding: 10px;
    height: calc(100% - 20px);
}

.vue-codemirror{
    height: 100%;
    text-align: left;
}

.CodeMirror{
    height: auto;
}

</style>
