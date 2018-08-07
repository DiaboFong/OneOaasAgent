<style lang="less">
    @import "./home.less";
    @import "../../styles/common.less";
</style>
<template>
    <div class="home-main">
        <Row :gutter="10">
            <Col :md="24" :lg="24">
                <Row :gutter="5">
                    <Col :xs="24" :sm="12" :md="6" :style="{marginBottom: '10px'}">
                        <infor-card
                            id-name="count_oday"
                            :end-val="count.today"
                            iconType="android-person-add"
                            color="#2d8cf0"
                            intro-text="agent今日安装"
                        ></infor-card>
                    </Col>
                    <Col :xs="24" :sm="12" :md="6" :style="{marginBottom: '10px'}">
                        <infor-card
                            id-name="count_count"
                            :end-val="count.count"
                            iconType="ios-eye"
                            color="#64d572"
                            :iconSize="50"
                            intro-text="agent安装总数"
                        ></infor-card>
                    </Col>
                    <Col :xs="24" :sm="12" :md="6" :style="{marginBottom: '10px'}">
                        <infor-card
                            id-name="count_connect"
                            :end-val="count.connect"
                            iconType="upload"
                            color="#ffd572"
                            intro-text="连接成功个数"
                        ></infor-card>
                    </Col>
                    <Col :xs="24" :sm="12" :md="6" :style="{marginBottom: '10px'}">
                        <infor-card
                            id-name="count_fail"
                            :end-val="count.fail"
                            iconType="shuffle"
                            color="#f25e43"
                            intro-text="连接异常个数"
                        ></infor-card>
                    </Col>
                </Row>
            </Col>
            <Col :md="24" :lg="24">
                <Row class-name="home-page-row1" :gutter="10">
                    <Col :md="12" :lg="24" :style="{marginBottom: '10px'}">
                        <Card>
                            <Row type="flex" class="user-infor">
                                <Col span="8">
                                    <Row class-name="made-child-con-middle" type="flex" align="middle">
                                        <img class="avator-img" src="../../images/avator.jpg" />
                                    </Row>
                                </Col>
                                <Col span="16" style="padding-left:6px;">
                                    <Row class-name="made-child-con-middle" type="flex" align="middle">
                                        <div>
                                            <b class="card-user-infor-name">{{ userinfo.user }}</b>
                                            <p>登录用户</p>
                                        </div>
                                    </Row>
                                </Col>
                            </Row>
                            <div class="line-gray"></div>
                            <Row class="margin-top-8">
                                <Col span="8"><p class="notwrap">上次登录时间:</p></Col>
                                <Col span="16" class="padding-left-8">{{ userinfo.login_time }}</Col>
                            </Row>
                            <Row class="margin-top-8">
                                <Col span="8"><p class="notwrap">上次登录IP:</p></Col>
                                <Col span="16" class="padding-left-8">{{ userinfo.login_ip }}</Col>
                            </Row>
                        </Card>
                    </Col>

                </Row>
            </Col>
        </Row>
    </div>
</template>

<script>
import Vue from "vue"
import cityData from './map-data/get-city-value.js';
import homeMap from './components/map.vue';
import dataSourcePie from './components/dataSourcePie.vue';
import visiteVolume from './components/visiteVolume.vue';
import serviceRequests from './components/serviceRequests.vue';
import userFlow from './components/userFlow.vue';
import countUp from './components/countUp.vue';
import inforCard from './components/inforCard.vue';
import mapDataTable from './components/mapDataTable.vue';
import toDoListItem from './components/toDoListItem.vue';

import VueResource from "vue-resource"
Vue.use(VueResource);
Vue.http.options.emulateJSON = true;
import Cookies from 'js-cookie';

export default {
    name: 'home',
    components: {
        homeMap,
        dataSourcePie,
        visiteVolume,
        serviceRequests,
        userFlow,
        countUp,
        inforCard,
        mapDataTable,
        toDoListItem
    },
    data () {
        return {
            userinfo: this.getUserInfo(),
            count: this.getCount(), //{connect: 0,count: 0,fail: 0,today: 0},
            cityData: cityData,
            showAddNewTodo: false,
            newToDoItemValue: ''
        };
    },
    computed: {
        avatorPath () {
            return localStorage.avatorImgPath;
        }
    },
    methods: {
        addNewToDoItem () {
            this.showAddNewTodo = true;
        },
        addNew () {
            if (this.newToDoItemValue.length !== 0) {
                this.toDoList.unshift({
                    title: this.newToDoItemValue
                });
                setTimeout(() => {
                    this.newToDoItemValue = '';
                }, 200);
                this.showAddNewTodo = false;
            } else {
                this.$Message.error('请输入待办事项内容');
            }
        },
        cancelAdd () {
            this.showAddNewTodo = false;
            this.newToDoItemValue = '';
        },
        getCount () {
            let ret = {
                "connect":0,
                "count":0,
                "fail":0,
                "today":0
            };
            Vue.http.get('/api/v1/dashboard/count')  
            .then(response => response.json())
            .then(response => {
                ret.connect =response.connect;
                ret.count =response.count;
                ret.fail =response.fail;
                ret.today =response.today;
            }).catch(err => {
                console.log('数据请求失败');
            });
            return ret;
        },
        getUserInfo() {
           let ret= {
                "user":"admin",
                "login_time":"2018-06-10 11:23:25",
                "login_ip":"127.0.0.1",
                "sessionid":""
            };
            if (Cookies.get('user') != 'undefined' ) {
                ret.user       = Cookies.get('user');
            }
            if (Cookies.get('login_ip') != 'undefined' ) {
                ret.login_ip   = Cookies.get('login_ip');
            }
            if  ( Cookies.get('login_time') != 'undefined') {
                ret.login_time = this.formatDate(Cookies.get('login_time'));
            }
            ret.sessionid  = Cookies.get('sessionid');
            
            return ret;
            
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
    }
};
</script>
