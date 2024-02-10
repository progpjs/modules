import React, {useContext, useEffect, useState, createContext} from "react";

function TestCore1() {
    function CompTheChild() {
        return <div>hello from the child!</div>;
    }

    function CompWithChildren({children}: any) {
        return <div><div>I have children</div>{children}</div>;
    }

    return <div className="myClass">
        <CompWithChildren><CompTheChild/></CompWithChildren>
    </div>
}

function TestCore2() {
    function increase(i: number): number {
        return i+1
    }

    function MyComponent({name}: any) {
        return <div>{name}</div>
    }

    let i = 0;
    const className = ["component-button"];

    return <div className={className.join(" ").trim()}>
        <MyComponent name="jojo" />
        <>Test #</>
        {increase(i)}
    </div>
}

function TestCore3() {
    const user = {
        name: 'Jo Joo',
        imageUrl: 'https://jo.jojo/yXOvdOSs.jpg',
        imageSize: 90,
    };

    return <>
        <h1>{user.name}</h1>
        <img
            className="avatar"
            src={user.imageUrl}
            alt={'Photo of ' + user.name}
            style={{
                width: user.imageSize,
                height: user.imageSize,
                backgroundImage: "url()"
            }}
        />
    </>
}

function TestStyle() {
    const style = {
        backgroundColor: "#ffffff",
        color: "red"
    }

    return <div style={style}></div>
}

function TestUseEffect() {
    // Is ignored when server side.
    useEffect(() => { console.log("effect are ignored server side") });

    return <div>TestStateEffect</div>
}

function TestUseState() {
    const [state, setState] = useState("my value");
    const [lazyState, setLazyState] = useState(() => "my lazy state");

    return <div>
        <div>State: {state}</div>
        <div>LazyState: {lazyState}</div>
    </div>
}

function TestContext() {
    let ThemeContext = createContext<string>("my default value");

    function UseThemeContext() {
        let theme = useContext(ThemeContext);
        return <div>{theme}</div>;
    }

    return <ThemeContext.Provider value="white">
        <ThemeContext.Provider value="pink">
            <UseThemeContext/>
        </ThemeContext.Provider>
    </ThemeContext.Provider>;
}

console.log(TestCore1().toString());
//console.log(TestCore2().toString());
//console.log(TestCore3().toString());
//console.log(TestStyle().toString());
//console.log(TestUseEffect().toString());
//console.log(TestUseState().toString());
//console.log(TestContext().toString());