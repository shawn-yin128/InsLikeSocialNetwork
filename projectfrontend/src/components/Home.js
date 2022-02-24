import React, { useState, useEffect } from "react";
import { Tabs, message, Row, Col, Button } from "antd";
import axios from "axios";

import SearchBar from "./SearchBar";
import PhotoGallery from "./PhotoGallery";
import { SEARCH_KEY, BASE_URL, TOKEN_KEY } from "../constants";

const { TabPane } = Tabs;

function Home(props) {
    const [posts, setPost] = useState([]);
    const [activeTab, setActiveTab] = useState("image");
    const [searchOption, setSearchOption] = useState({
        type: SEARCH_KEY.all,
        keyword: ""
    });

    // when searchOption changes, call the arrow function
    useEffect(() => {
        const { type, keyword } = searchOption;
        fetchPost(searchOption);
    }, [searchOption]);

    const fetchPost = (option) => {
        const { type, keyword } = option;
        let url = "";

        if (type === SEARCH_KEY.all) {
            url = `${BASE_URL}/search`;
        } else if (type === SEARCH_KEY.user) {
            url = `${BASE_URL}/search?user=${keyword}`;
        } else {
            url = `${BASE_URL}/search?keywords=${keyword}`;
        }

        const opt = {
            method: "GET",
            url: url,
            headers: {
                Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`
            }
        };

        axios(opt)
            .then((res) => {
            if (res.status === 200) {
                setPost(res.data);
            }
            })
            .catch((err) => {
            message.error("Fetch posts failed!");
            console.log("fetch posts failed: ", err.message);
            });
    };

    const renderPosts = (type) => {
        if (!posts || posts.length === 0) {
            return <div>No data!</div>;
        }
        if (type === "image") {
            console.log("images -> ", posts);
            return "images";
        } else if (type === "video") {
            console.log("video -> ", posts);
            return "videos";
        }
    };

    const operations = <Button>Upload</Button>;
    return (
        <div className="home">
            <SearchBar />
            <div className="display">
            <Tabs
                onChange={(key) => setActiveTab(key)}
                defaultActiveKey="image"
                activeKey={activeTab}
                tabBarExtraContent={operations}
            >
                <TabPane tab="Images" key="image">
                {renderPosts("image")}
                </TabPane>
                <TabPane tab="Videos" key="video">
                {renderPosts("video")}
                </TabPane>
            </Tabs>
            </div>
        </div>
    );
    }

export default Home;