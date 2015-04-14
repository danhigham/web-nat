var Template = function() {
    this.cached = {};
};

var T = new Template();
$.extend(Template.prototype, {
    render: function(name, callback) {
        if (T.isCached(name)) {
            callback(T.cached[name]);
        } else {
            $.get(T.urlFor(name), function(raw) {
                T.store(name, raw);
                T.render(name, callback);
            });
        }
    },
    renderSync: function(name, callback) {
        if (!T.isCached(name)) {
            T.fetch(name);
        }
        T.render(name, callback);
    },
    prefetch: function(name) {
        $.get(T.urlFor(name), function(raw) { 
            T.store(name, raw);
        });
    },
    fetch: function(name) {
        if (! T.isCached(name)) {
            var raw = $.ajax({'url': T.urlFor(name), 'async': false}).responseText;
            T.store(name, raw);         
        }
    },
    isCached: function(name) {
        return !!T.cached[name];
    },
    store: function(name, raw) {
        T.cached[name] = Handlebars.compile(raw);
    },
    urlFor: function(name) {
        return "/assets/templates/"+ name + ".handlebars";
    }
});

var svg = d3.select("svg"),
    inner = svg.select("g"),
    zoom = d3.behavior.zoom().on("zoom", function() {
        inner.attr("transform", "translate(" + d3.event.translate + ")" +
            "scale(" + d3.event.scale + ")");
    });
svg.call(zoom);

var render = new dagreD3.render();

var g = new dagreD3.graphlib.Graph();
g.setGraph({
    nodesep: 70,
    ranksep: 50,
    rankdir: "LR",
    marginx: 20,
    marginy: 20
});

var containers = [
{ id: "Wordpress 1234", host: "host 1", port_nat: [ { container_port: 80, port: 8080, from: "192.168.3.0/24"} ] },
{ id: "Minecraft Server", host: "host 2", port_nat: [ { container_port: 1234, port: 1234, from: "anywhere"} ] },
{ id: "MySQL", host: "host 1", port_nat: [ { container_port: 3306, port: 3306, from: "192.168.3.0/24"} ] }
];

function htmlNode(caption, subtext, className) {
    var res;

    context = {
        caption: caption,
        subtext: subtext
    }
    T.renderSync('node', function(t) {
        res = t(context);
    });

    return res;
}

function draw(isUpdate) {
    // add NAT box
    //
    var html = htmlNode("NAT", "NAT vm", "nat");
    console.log(html);
    g.setNode("nat-box", {
        labelType: "html",
        label: html,
        rx: 5,
        ry: 5,
        padding: 0
        //class: "nat-box"
    });

    for (var id in containers) {
        var container = containers[id]        
        var html = htmlNode(container.id, container.id, "container");

        g.setNode(id, {
            labelType: "html",
            label: html,
            rx: 5,
            ry: 5,
            padding: 0
            // class: className
        });

        if (container.port_nat) {
            _.each(container.port_nat, function(nat) {
                g.setEdge("nat-box", id, {
                    label: nat.port,
                    width: 40
                });
            });
        }
    }

    inner.call(render, g);

    var zoomScale = zoom.scale();
    var graphWidth = g.graph().width + 80;
    var graphHeight = g.graph().height + 40;
    var width = parseInt(svg.style("width").replace(/px/, ""));
    var height = parseInt(svg.style("height").replace(/px/, ""));
    zoomScale = Math.min(width / graphWidth, height / graphHeight);
    var translate = [(width/2) - ((graphWidth*zoomScale)/2), (height/2) - ((graphHeight*zoomScale)/2)];
    zoom.translate(translate);
    zoom.scale(zoomScale);
    zoom.event(isUpdate ? svg.transition().duration(500) : d3.select("svg"));
}

draw(false);
