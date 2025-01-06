// app.js
class DiaryApp {
    constructor() {
        this.initLightBackground();
        this.initDraggableModules();
        this.initInspirationThemes();
        this.initGrowthTracking();
    }

    // 动态光影效果
    initLightBackground() {
        const background = document.getElementById('light-background');
        const hour = new Date().getHours();

        const lightStyles = {
            morning: 'linear-gradient(45deg, #FFD700, #87CEEB)',
            afternoon: 'linear-gradient(45deg, #FFA500, #FF6347)',
            evening: 'linear-gradient(45deg, #4B0082, #191970)'
        };

        if (hour >= 5 && hour < 12) {
            background.style.background = lightStyles.morning;
        } else if (hour >= 12 && hour < 18) {
            background.style.background = lightStyles.afternoon;
        } else {
            background.style.background = lightStyles.evening;
        }
    }

    // 模块拖拽功能
    initDraggableModules() {
        const modules = document.querySelectorAll('.diary-module');

        modules.forEach(module => {
            module.addEventListener('dragstart', this.handleDragStart);
            module.addEventListener('dragend', this.handleDragEnd);
        });
    }

    // 灵感主题推送
    initInspirationThemes() {
        const themes = [
            '春日花事',
            '旅行足迹',
            '职场成长',
            '生日感悟'
        ];

        const themeList = document.querySelector('.theme-list');
        themes.forEach(theme => {
            const themeItem = document.createElement('div');
            themeItem.textContent = theme;
            themeList.appendChild(themeItem);
        });
    }

    // 成长轨迹可视化
    initGrowthTracking() {
        // 这里可以接入数据可视化库，如 ECharts
        const timeline = document.querySelector('.timeline');
        const emotionChart = document.querySelector('.emotion-chart');

        // 模拟数据展示
        timeline.innerHTML = '成长时间轴';
        emotionChart.innerHTML = '情绪变化图表';
    }

    handleDragStart(e) {
        e.target.style.opacity = '0.5';
    }

    handleDragEnd(e) {
        e.target.style.opacity = '1';
    }
}

// 初始化应用
document.addEventListener('DOMContentLoaded', () => {
    new DiaryApp();
});