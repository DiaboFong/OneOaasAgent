<template>
  <div>
         <Row>
            <Col span="24">
                <Card>
                    <p slot="title">
                        <Icon type="ios-cloud-upload-outline"></Icon>
                        基本上传功能
                        <Icon type="ios-cloud-upload-outline">
                        <Upload
                            multiple
                            action="//jsonplaceholder.typicode.com/posts/">
                            <Button type="ghost" icon="ios-cloud-upload-outline">Upload files</Button>
                        </Upload>
                        </Icon>
                        

                    </p>
                    <div class="height-120px">
                            <Card>
                                    <can-edit-table refs="table1" @on-delete="handleDel" v-model="tableData" :columns-list="columnsList"></can-edit-table>
                            </Card>
                    </div>
                </Card>
            </Col>
          </Row>
  </div>
</template>

<script>
import canEditTable from './components/canEditTable.vue';
import tableData from './components/table_data.js';
export default {
    name: 'editable-table',
    components: {
        canEditTable
    },
    data () {
        return {
            columnsList: [],
            tableData: [],
            editInlineColumns: [],
            editInlineData: [],
            editIncellColumns: [],
            editIncellData: [],
            editInlineAndCellColumn: [],
            editInlineAndCellData: [],
            showCurrentColumns: [],
            showCurrentTableData: false
        };
    },
    methods: {
        getData () {
            this.columnsList = tableData.table1Columns;
            this.tableData = tableData.table1Data;
            this.editInlineColumns = tableData.editInlineColumns;
            this.editInlineData = tableData.editInlineData;
            this.editIncellColumns = tableData.editIncellColumns;
            this.editIncellData = tableData.editIncellData;
            this.editInlineAndCellColumn = tableData.editInlineAndCellColumn;
            this.editInlineAndCellData = tableData.editInlineAndCellData;
            this.showCurrentColumns = tableData.showCurrentColumns;
        },
        handleNetConnect (state) {
            this.breakConnect = state;
        },
        handleLowSpeed (state) {
            this.lowNetSpeed = state;
        },
        getCurrentData () {
            this.showCurrentTableData = true;
        },
        handleDel (val, index) {
            this.$Message.success('删除了第' + (index + 1) + '行数据');
        },
        handleCellChange (val, index, key) {
            this.$Message.success('修改了第 ' + (index + 1) + ' 行列名为 ' + key + ' 的数据');
        },
        handleChange (val, index) {
            this.$Message.success('修改了第' + (index + 1) + '行数据');
        }
    },
    created () {
        this.getData();
    }
};
</script>

<style>
  
</style>
