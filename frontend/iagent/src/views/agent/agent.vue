<template>
  <div>
    <Card shadow>
       <Form ref="formInstall" :model="formInstall" :rules="ruleInstall" inline>
            <FormItem prop="installip">
                <Input style="width: 130px" v-model="formInstall.installip" type="textarea" :autosize="{minRows: 1,maxRows: 5}" placeholder="输入IP地址"></Input>
            </FormItem>
            <FormItem prop="user">
                <Input style="width: 100px" v-model="formInstall.user" placeholder="请输入用户名"></Input>
            </FormItem>
            <FormItem prop="password">
                <Input style="width: 100px" v-model="formInstall.password" type="password" placeholder="请输入密码"></Input>
            </FormItem>
            <!--Select style="width:90px" value="v1">
                <Option v-for="item in scriptsList" :value="item.value" :key="item.value">{{ item.value }}</Option>
            </Select-->
            <FormItem>
                <Button @click="install"  v-show="installLoading" icon="ios-checkmark" style="width:90px;" type="primary">安装</Button>
            </FormItem>
               <Button style="float: right;" icon="ios-search" @click="search" v-show="searchLoading" >搜　索</Button>
               <Button class="search-button" v-show="!searchLoading">
                <Icon type="load-c" class="demo-spin-icon-load"></Icon>
                <span>Loading</span>
               </Button>
               <Input  style="width: 130px;float: right;" v-model="searchip" placeholder="请输入IP" @keyup.13.native="search" icon="ios-search"></Input>
        </Form>
        <Modal title="安装日志" @scrollable="true" v-model="modal2" width="950px"  style="{height: rowHeight + 'px'}">
            <p slot="header">安装日志</p>
            <Tabs> 
                <Input  readonly="true" v-model="agentLog" type="textarea" :autosize="{minRows: 1,maxRows: 20}" placeholder="agent安装日志"></Input>
   
            </Tabs>
            <div slot="footer">
                <Button type="primary" size="large" @click="modal2=!modal2">关闭</Button>
            </div>
        </Modal>
        
    
        <Table :data="tableData1" :columns="tableColumns1" stripe id="resultDatatables"></Table>
        <div style="margin: 10px;overflow: hidden">
            <div style="float: right;">
                <Page :total="dataListLength" :current="1" @on-change="changePage" show-sizer show-total @on-page-size-change="changePageSize" :page-size="pageSize" placement="top"></Page>
            </div>
        </div>
  
  
  
    </Card>
  </div>
</template>

<script>
    import Vue from "vue"
    import VueResource from "vue-resource"
    let Base64 = require('js-base64').Base64;
    Vue.use(VueResource);
    Vue.http.options.emulateJSON = true;
    import { Modal } from 'iview'
    import { iviewRow, iviewCol, iviewButton, iviewInput, iviewForm, iviewFormItem } from 'iview'
    import Cookies from 'js-cookie';
    //import moment from 'moment'

    export default {
        data () {
            return {
                searchip:'',           
                formInstall: {
                    installip: '',
                    user:'',
                    password:'',
                   
                },
                ruleInstall: {
                    installip: [
                        { required: true, message: '请输入IP地址', trigger: 'blur' }
                    ],
                    user: [
                        { required: true, message: '请输入用户名', trigger: 'blur' }
                    ],
                    password: [
                        { required: true, message: '请输入密码', trigger: 'blur' },
                        { type: 'string', min: 1, message: '密码少于1位', trigger: 'blur' }
                    ]
                },               
                scriptsList: [{value: 'v1'}, {value: 'v2'}],
                agentLog:'',
                modal2: false,
                modal1: true,
                resultData: [],
                pageData: [],
                dataListLength: 0,
                pageSize: 10,
                searchip: '',
                showData: false,
                searchLoading: true,
                installLoading: true,
                syncLoading: false,
                tableWidth: false,
                // 表格
                tableData1: this.getAgentLog(),
                tableColumns1: [
                    {
                        title: '安装时间',
                        key: 'update_time',
                        fixed: 'left',
                        //width:180,
                        sortable: true,
                        sortType: "desc",
                        render: (h, params) => {
                             return h('div', this.tableData1[params.index].update_time);
                        }
                    },
                    {
                        title: 'IP地址',
                        key: 'ip',
                        fixed: 'left',
                        //width:110,
                        sortable: true,
                        render: (h, params) => {
                            return h('div', this.tableData1[params.index].ip);
                        }
                    },
                    
                    {
                        title: '操作详情',
                        key: 'status',
                        fixed: 'left',
                        //width: 150,
                        sortable: true,
                        render: (h, params) => {
                            var _type="";
                            var _icon="";
                            if (params.row.status === "认证失败") {
                                _type = "error"
                                _icon = "heart-broken"
                            } else if (params.row.status === "脚本已执行") {
                                _type = "info"
                                _icon = "coffee"
                            } else if (params.row.status === "无法连接") {
                                _type = "disabled"
                                _icon = "android-alert"
                            };
                             return h('div', [
                                h('Button', {
                                    props: {
                                        icon: _icon,
                                        type: _type,
                                        size: 'small'
                                    },
                                    style: {
                                        marginRight: '10px'
                                    }
                                }, params.row.status )
                             ]);
                        }
                        /*render: (h, params) => {
                            //认证失败
                            //无法连接
                            //脚本已执行
                            return h('div', this.tableData1[params.index].status);
                        }*/
                    },
                    {
                        title: 'agent状态',
                        key: 'connect',
                        fixed: 'left',
                        //width: 150,
                        sortable: true,
                        render: (h, params) => {
                            var _type="";
                            var _icon="";
                            if (params.row.connect === "未知") {
                                _type = "dashed"
                                _icon = "help-circled"
                            } else if (params.row.connect === "连接") {
                                _type = "success"
                                _icon = "checkmark-circled"
                            } else if (params.row.connect === "离线") {
                                _type = "warning"
                                _icon = "close-circled"
                            };
                             return h('div', [
                                h('Button', {
                                    props: {
                                        icon: _icon,
                                        type: _type,
                                        size: 'small'
                                    },
                                    style: {
                                        marginRight: '10px'
                                    }
                                }, params.row.connect )
                             ]);
                        }
                    },
                    {
                        title: '操作',
                        key: 'action',
                        align: 'center',
                        render: (h, params) => {
                            return h('div', [
                                h('Button', {
                                    props: {
                                        icon: "eye",
                                        type: 'primary',
                                        size: 'small'
                                    },
                                    style: {
                                        marginRight: '10px'
                                    },
                                    on: {
                                        click: () => {
                                            this.showModal(params.index)
                                        }
                                    }
                                }, '查看'),
                                h('Button', {
                                    shape: "circle",
                                    props: {
                                        icon: "android-delete",
                                        type: 'error',
                                        size: 'small'
                                    },
                                    style: {
                                        marginRight: '5px'
                                    },
                                    on: {
                                        click: () => {
                                            this.remove(params.index)
                                        }
                                    }
                                }, '卸载')
                            ]);
                        }
                    }
                ],
            }
        },
        computed: {
            rowHeight: function () {
                return this.$store.state.interfaceConfig.appHeight / 2 - 20
            },
            modalWidth: function () {
                return this.$store.state.interfaceConfig.appWidth/2-100
            }

        },
        methods: {
            start () {
                this.$Loading.start();
            },
            finish () {
                this.$Loading.finish();
            },
            error () {
                this.$Loading.error();
            },
            msg() {
                return "iagent"
            },
            sync(index) {
                this.syncLoading = true;
                let ipaddress = this.tableData1[index].ipaddress;
                const urlsync = "/api/log/syncfromstore?store="+ipaddress
                this.start();
                Vue.http.get(urlsync)
                .then(response => response.json())
                .then(response => {
                    this.finish();
                }).catch(err => {
                    console.log('同步请求失败：'+err.status+','+err.statusText);
                });
                this.syncLoading = false;
            },
            showModal (index) {
                this.modal2 = true;
                let ip = this.tableData1[index].ip;
 
                // 日志上下文数据
                const url = "/api/v1/agent/logbyip?ip=" + ip
                Vue.http.get(url)
                .then(response => response.json())
                .then(response => {
                        this.agentLog=response.content.log;
                }).catch(err => {
                        console.log('查询日志失败：'+err.status+','+err.statusText);
                });
                //this.modal = true
            },
            getAgentLog () {
                //this.showData = true;
                let data1 = [];
                this.start();
            
                Vue.http.get('/api/v1/agent/log')  
                .then(response => response.json())
                .then(response => {
                    let data = response.content;
                    let len = this.pageSize;
                    if(data.length <= this.pageSize){
                        len = data.length;
                    }
                    for(let i = 0; i < len; i++) {
                        data1.push({
                            "ip":data[i].ip,
                            "update_time":this.formatDate(data[i].update_time),
                            "status": data[i].status,
                            "connect": data[i].connect
                        })
                    }
                    this.dataListLength = data.length;
                    this.resultData = data;
                }).catch(err => {
                        console.log('数据请求失败：'+err.status+','+err.statusText);
                });
                this.finish();
                //this.showData = false; 
                //
                /*data1 = [{
                        "ipaddress": "192.168.4.56",
                        "installtime": "2018-05-23 11:23:05",
                        "installdetail": "安装成功",
                        "agentstatus": "已安装"
                    },{
                        "ipaddress": "192.168.4.57",
                        "installtime": "2018-05-23 11:23:07",
                        "installdetail": "安装成功",
                        "agentstatus": "已安装"
                    },{
                        "ipaddress": "192.168.4.58",
                        "installtime": "2018-05-23 11:23:08",
                        "installdetail": "安装成功",
                        "agentstatus": "已安装"
                    }]
                
                
                */
                return data1;
            },
            remove (index) {
                const urlremove = "/api/v1/agent/remove";
                var params = {
                    "ip":this.tableData1[index].ip
                };
                this.$Modal.confirm({
                title: "提示",
                content: "<p>确定卸载"+params.ip+"吗？</p>",
                onOk: () => {
                    Vue.http.post(urlremove,params)
                    .then(response => response.json())
                    .then(response => {
                            if (response.code === 200 ) {
                                 this.$Message.success('执行:\t'+response.msg);
                            }
                            if (response.code !== 200 ) {
                                 this.$Message.error('失败:\t'+response.msg);
                            }
                        //this.data.splice(index, 1);
                        //this.$Message.success("卸载成功");
                    }).catch(err => {
                        this.$Message.error('卸载请求失败');
                    });
                    var i = 0;
                    var IntervalRemove = setInterval(() => {
                        this.tableData1=this.getAgentLog();
                        i++;
                        if (i > 2) {
                            //删除定时器IntervalName
                            clearInterval(IntervalRemove);
                        };
                    }, 5000); 
                },
                onCancel: () => {}
                
            });

            },
            formatDate (datestr) {
                var date= new Date(datestr*1000);
                const y = date.getFullYear();
                let m = date.getMonth() + 1;
                m = m < 10 ? '0' + m : m;
                let d = date.getDate();
                d = d < 10 ? ('0' + d) : d;
                let h = date.getHours();
                h = h < 10 ? ('0' + h) : h;
                let M = date.getMinutes();
                M = M < 10 ? ('0' + M) : M;
                let s = date.getSeconds();
                s = s < 10 ? ('0' + s) : s;
                return y + '-' + m + '-' + d+' '+h+':'+M+':'+s;               
            },
            changePageSize (size){
                this.pageSize = size;
                this.changePage();
            },
            changePage (page) {
                if(!page){
                    page = 1;
                }
                let data1 = [];
                var _start = ( page - 1 ) * this.pageSize;
                var _end = page * this.pageSize;
                var data = this.resultData;
                if(data.length<=_end){
                    _end = data.length;
                }
                for(let i = _start; i < _end; i++) {
                    data1.push({
                        "ip":data[i].ip,
                        "update_time":this.formatDate(data[i].update_time),
                        "status": data[i].status,
                        "connect": data[i].connect
                    })
                }
                this.tableData1 = data1;
            },
            search () {
                console.log(this.searchip);
                if (this.searchip == "" ) {
                    this.$Message.error('请输入IP地址....');
                    return;
                };
                this.searchLoading = false;
                let data1 = [];
                this.start();
                const searchurl = "/api/v1/agent/logbyip?ip="+this.searchip;
                Vue.http.get(searchurl)
                .then(response => response.json())
                .then(response => {
                    let data=response.content;
                    /*let len = this.pageSize;
                    if(data.id === 0){
                        this.dataListLength = 0;
                        this.tableData1 = data1;
                        return;
                    }
                    if(data.length <= this.pageSize){
                        len = data.length;
                    }
                    for(let i = 0; i < len; i++) {
                        data1.push({
                            "ip":data[i].ip,
                            "update_time":data[i].update_time,
                            "status": data[i].status,
                            "connect": data[i].connect
                        });
                    }*/
                    //this.dataListLength = data.length;
                    //this.resultData = data;
                    data1.push({
                            "ip":data.ip,
                            "update_time":this.formatDate(data[i].update_time),
                            "status": data.status,
                            "connect": data.connect
                        })
                    this.tableData1 = data1;
                }).catch(err => {
                        this.$Message.error('搜索'+this.searchip+'请求失败......\t'+err.status+','+err.statusText);
                        //console.log('搜索请求失败：'+err.status+','+err.statusText);
                });
                this.searchLoading = true;
                this.finish();
            },
            install () {
                this.$refs.formInstall.validate((valid) => {
                    if (valid) {                       
                        this.installLoading = false;
                        let data1 = [];
                        this.start();
                        const installipaddress = "/api/v1/agent/install"; //"?installtime=this.installip";
                        var params = {
                            "ip": this.formInstall.installip.split("\n").join(),
                            "username": Base64.encode(this.formInstall.user),
                            "password": Base64.encode(this.formInstall.password)
                        }; 
                        Vue.http.post(installipaddress,params)
                        .then(response => response.json())
                        .then(response => {
                            if (response.code === 200 ) {
                                 this.$Message.success('执行:\t'+response.msg);
                                
                            }
                            if (response.code !== 200 ) {
                                 this.$Message.error('失败:\t'+response.msg);
                            }
                            
                        /*let data=response.data;
                        let len = this.pageSize;
                        if(data.id === 0){
                            this.dataListLength = 0;
                            this.tableData1 = data1;
                            return;
                        }
                        if(data.length <= this.pageSize){
                            len = data.length;
                        }
                        for(let i = 0; i < len; i++) {
                            data1.push({
                                "ipaddress":data[i].ipaddress,
                                "installtime":data[i].installtime,
                                "installdetail": data[i].installdetail,
                                "agentstatus": data[i].agentstatus
                            });
                        }
                        this.dataListLength = data.length;
                        this.resultData = data;
                        this.tableData1 = data1;*/
                        }).catch(err => {
                                this.$Message.error('请求数据失败....');
                                console.log('安装请求失败：'+err.status+','+err.statusText);
                        });
                        this.installLoading = true;
                        //this.$Message.success('提交安装信息成功!');
                        //window.location.reload();
                        this.finish();
                        //setInterval(this.getAgentLog, 3000);
                        var i = 0;
                        var IntervalGetLog = setInterval(() => {
                            //this.getAgentLog;
                            //console.log(this.getAgentLog());
                            this.tableData1=this.getAgentLog();
                            i++;
                            if (i > 5) {
                                //删除定时器IntervalName
                                clearInterval(IntervalGetLog);
                            };
                        }, 10000);    
                    } else {
                        this.$Message.error('提交数据格式错误!');
                    }
                });


            },
        },
        mounted() {
            const _this = this;
            window.onresize = function temp() {
                let offsetWidth = `${document.body.offsetWidth}`;
                if( offsetWidth <= 1500 ){
                    _this.tableWidth = true;
                }else{
                    _this.tableWidth = false;
                }
            };
        },
        created() {
            Cookies.set('install', '0');
        },
    }
</script>
<style scoped>
    /* 布局 */
    .layout-breadcrumb{
        padding: 10px 15px 0;
    }
    .layout-content{
        min-height: 800px;
        margin: 15px;
        overflow: hidden;
        background: #fff;
        border-radius: 4px;
    }
    .layout-content-main{
        padding: 10px;
    }
    .layout-ceiling-main a{
        color: #9ba7b5;
    }
    /* loading */
    .demo-spin-icon-load{
        animation: ani-demo-spin 1s linear infinite;
    }
    @keyframes ani-demo-spin {
        from { transform: rotate(0deg);}
        50%  { transform: rotate(180deg);}
        to   { transform: rotate(360deg);}
    }
    .demo-spin-col{
        height: 100px;
        position: relative;
        border: 1px solid #eee;
    }
    .loading-data{
        position: fixed;
        z-index: 100;
        width: 100%;
        height: 100vh;
        top: 0;left: 0;
    }
    .search-style{
        float: right;
        margin-right: 26px;
    }

    .search-button{
        width:85px;
    }
    /* button */
    .button-circle{
        height: 102px;
        border-radius: 50%;
        margin-bottom: 40px;
    }
    /* 表格 */
    .table-style{
        text-align: left;
        font-size: 14px;
        width: 100%;
        border-collapse:collapse;
    }
    .table-style th{
        height:40px;
        border: 1px solid #e5e5e5;
    }
    .table-style td{
        width: 28%;
        height: 40px;
        padding: 0 30px;
        border: 1px solid #e5e5e5;
        overflow: auto;
        color: #666;
    }
    .table-style td:first-child{
        width:16%;
        color:#333;
    }
    /* 门店日志图片大小 */
    .store-width{
        width: 80px;
        padding-top: 10px;
    }
    .store-log{
        font-size: 14px;
        position: absolute;
        left: 85px;
        top: 85px;
        color: #2180ff;
    }
    /* 日志错误 */
    .store-error-log{
        position: absolute;
        top: -50px;
    }
    /* 数据不同高亮显示 */
    .diff-data{
        background: rgba(33, 128, 255, 0.1);
    }
</style>
